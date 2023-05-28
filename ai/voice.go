package ai

type VoiceSelector struct {
	voicesMap    map[string]map[int]string
	defaultVoice string
}

func NewVoiceSelector(voicesMap map[string]map[int]string, defaultVoice string) *VoiceSelector {
	return &VoiceSelector{
		voicesMap:    voicesMap,
		defaultVoice: defaultVoice,
	}
}

func (v *VoiceSelector) GetVoice(language string, id int) (string, error) {
	voice, ok := v.voicesMap[language][id]
	if !ok {
		return v.defaultVoice, nil
	}

	return voice, nil

}
