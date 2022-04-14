package utils

import "strings"

var (
	record_static = true

	// http static resource file extension
	static_ext = []string{
		"js",
		"css",
		"ico",
	}

	// media resource files type
	media_types = []string{
		"image",
		"video",
		"audio",
	}

	// http static resource files
	static_types = []string{
		"text/css",
		// "application/javascript",
		// "application/x-javascript",
		"application/msword",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/x-ms-wmd",
		"application/x-shockwave-flash",
	}
)

func IsStatic(contentType, extension string) bool {
	var mtype string
	if contentType != "" {
		mtype = strings.Split(contentType, "/")[0]
	}
	if ContainsString(static_ext, extension) {
		return true
	} else if ContainsString(static_types, contentType) {
		return true
	} else if ContainsString(media_types, mtype) {
		return true
	}
	return false
}

func ContainsString(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
