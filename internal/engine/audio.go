package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 44100

// audioManager handles sound loading and playback
type AudioManager struct {
	context *audio.Context
	sounds  map[string]*SoundAsset
	currentMusic string
}

// soundAsset represents a loaded sound
type SoundAsset struct {
	Name     string
	Volume   float64
	Pan      float64
	Player   *audio.Player
	FilePath string
	Looping  bool
}

func NewAudioManager() *AudioManager {
	return &AudioManager{
		context: audio.NewContext(sampleRate),
		sounds:  make(map[string]*SoundAsset),
	}
}

// loadSoundsFromManifest loads sound metadata and prepares for playback
func (am *AudioManager) LoadSoundsFromManifest(assets *AssetManager) error {
	metas := assets.GetSoundMeta()
	basePath := assets.GetBasePath()

	loaded := 0
	for _, meta := range metas {
		if meta.FileName == "" {
			continue
		}

		soundPath := filepath.Join(basePath, "sounds", meta.FileName)
		if _, err := os.Stat(soundPath); os.IsNotExist(err) {
			continue
		}

		am.sounds[meta.Name] = &SoundAsset{
			Name:     meta.Name,
			Volume:   meta.Volume,
			Pan:      meta.Pan,
			FilePath: soundPath,
		}
		loaded++
	}

	fmt.Printf("Registered %d sounds for playback\n", loaded)
	return nil
}

// setVolume sets the volume for a named sound (0.0 - 1.0)
func (am *AudioManager) SetVolume(name string, vol float64) {
	sa, ok := am.sounds[name]
	if !ok {
		return
	}
	sa.Volume = vol
	if sa.Player != nil {
		sa.Player.SetVolume(vol)
	}
}

// play plays a sound by name. If looping is true, it will restart when finished (best-effort).
func (am *AudioManager) Play(name string, looping ...bool) {
	sa, ok := am.sounds[name]
	if !ok {
		return
	}
	sa.Looping = len(looping) > 0 && looping[0]

	// stop any existing player for this sound to prevent double-play
	if sa.Player != nil && sa.Player.IsPlaying() {
		sa.Player.Pause()
		sa.Player = nil
	}

	// load and play the audio file
	f, err := os.Open(sa.FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: cannot open sound %s: %v\n", name, err)
		return
	}

	ext := strings.ToLower(filepath.Ext(sa.FilePath))

	switch ext {
	case ".wav":
		s, err := wav.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "WARNING: cannot decode wav %s: %v\n", name, err)
			return
		}
		player, err := am.context.NewPlayer(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: cannot create player for %s: %v\n", name, err)
			return
		}
		player.SetVolume(sa.Volume)
		player.Play()
		sa.Player = player
	case ".mp3":
		s, err := mp3.DecodeWithSampleRate(sampleRate, f)
		if err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "WARNING: cannot decode mp3 %s: %v\n", name, err)
			return
		}
		player, err := am.context.NewPlayer(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: cannot create player for %s: %v\n", name, err)
			return
		}
		player.SetVolume(sa.Volume)
		player.Play()
		sa.Player = player
	default:
		f.Close()
		fmt.Fprintf(os.Stderr, "WARNING: unsupported audio format %s for %s\n", ext, name)
		return
	}
}

// update should be called each frame to service looping sounds.
func (am *AudioManager) Update() {
	for _, sa := range am.sounds {
		if sa == nil || sa.Player == nil || !sa.Looping {
			continue
		}
		if !sa.Player.IsPlaying() {
			_ = sa.Player.Rewind()
			sa.Player.Play()
		}
	}
}

// playMusic plays exactly one looping BGM track at a time.
func (am *AudioManager) PlayMusic(name string) {
	if name == "" {
		return
	}
	if am.currentMusic == name && am.IsPlaying(name) {
		return
	}
	if am.currentMusic != "" && am.currentMusic != name {
		am.Stop(am.currentMusic)
	}
	am.Play(name, true)
	am.currentMusic = name
}

// stopMusic stops the active BGM track, if any.
func (am *AudioManager) StopMusic() {
	if am.currentMusic == "" {
		return
	}
	am.Stop(am.currentMusic)
	am.currentMusic = ""
}

// stop stops a playing sound
func (am *AudioManager) Stop(name string) {
	sa, ok := am.sounds[name]
	if !ok || sa.Player == nil {
		return
	}
	sa.Looping = false
	sa.Player.Pause()
}

// stopAll stops all playing sounds
func (am *AudioManager) StopAll() {
	for _, sa := range am.sounds {
		if sa.Player != nil {
			sa.Looping = false
			sa.Player.Pause()
		}
	}
	am.currentMusic = ""
}

// isPlaying checks if a sound is currently playing
func (am *AudioManager) IsPlaying(name string) bool {
	sa, ok := am.sounds[name]
	if !ok || sa.Player == nil {
		return false
	}
	return sa.Player.IsPlaying()
}

// hasSound reports whether a named sound is available in the manifest/cache.
func (am *AudioManager) HasSound(name string) bool {
	_, ok := am.sounds[name]
	return ok
}
