package fmtp

type hevcFMTP struct {
	parameters map[string]string
}

func (h *hevcFMTP) MimeType() string {
	return "video/h265"
}

// Match returns true if h and b are compatible fmtp descriptions
// Based on RFC6184 Section 8.2.2:
//   The parameters identifying a media format configuration for HEVC
//   are profile-level-id and packetization-mode.  These media format
//   configuration parameters (except for the level part of profile-
//   level-id) MUST be used symmetrically; that is, the answerer MUST
//   either maintain all configuration parameters or remove the media
//   format (payload type) completely if one or more of the parameter
//   values are not supported.
//     Informative note: The requirement for symmetric use does not
//     apply for the level part of profile-level-id and does not apply
//     for the other stream properties and capability parameters.
func (h *hevcFMTP) Match(b FMTP) bool {
	c, ok := b.(*hevcFMTP)
	if !ok {
		return false
	}

	// test packetization-mode
	hpmode, hok := h.parameters["packetization-mode"]
	if !hok {
		return false
	}
	cpmode, cok := c.parameters["packetization-mode"]
	if !cok {
		return false
	}

	if hpmode != cpmode {
		return false
	}

	// test profile-level-id
	hplid, hok := h.parameters["profile-level-id"]
	if !hok {
		return false
	}

	cplid, cok := c.parameters["profile-level-id"]
	if !cok {
		return false
	}

	if !profileLevelIDMatches(hplid, cplid) {
		return false
	}

	return true
}

func (h *hevcFMTP) Parameter(key string) (string, bool) {
	v, ok := h.parameters[key]
	return v, ok
}
