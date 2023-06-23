package ytdlp

import (
	"fmt"
	"testing"
)

func TestIsDownloadable(t *testing.T) {
	type testCase struct {
		URL            string
		ytType         YtdlpURLType
		isDownloadable bool
	}

	// TODO: Write tests for situation when you have URL in some text
	tests := []testCase{
		{"https://www.youtube.com/watch?v=123456789zxcvbnasdfqwew", YoutubeVideoType, true},
		{"www.youtube.com/watch?v=123456789zxcvbnasdfqwew", YoutubeVideoType, true},
		{"youtube.com/watch?v=123456789zxcvbnasdfqwew", YoutubeVideoType, true},
		{"https://youtu.be/312u8fdjaf", YoutubeVideoType, true},
		{"youtu.be/312u8fdjaf", YoutubeVideoType, true},

		{"https://www.youtube.com/playlist?list=1234345sdfjlka", YoutubePlaylistType, true},
		{"www.youtube.com/playlist?list=1234345sdfjlka", YoutubePlaylistType, true},
		{"youtube.com/playlist?list=1234345sdfjlka", YoutubePlaylistType, true},
		{"youtube.com/playlist?list=0", YoutubePlaylistType, true},

		{"https://vk.com/video-123891358787_3199312", VkVideoType, true},
		{"vk.com/video-123891358787_3199312", VkVideoType, true},
		{"vk.com/video-0", VkVideoType, true},
		{"vk.com/videos-", VkVideoType, true},

		{"youtube./watch?v=123456789zxcvbnasdfqwew", UndefinedType, false},
		{"youtube.com/playlist?list=", UndefinedType, false},
		{"vk.com/video", UndefinedType, false},
	}

	yt := NewDefault()
	for _, test := range tests {
		testname := fmt.Sprintf("Trying to download %s with %s type is  %v downloadable", test.URL, test.ytType, test.isDownloadable)
		t.Run(testname, func(t *testing.T) {
			tp, _, down := yt.IsDownloadable(test.URL)
			if down != test.isDownloadable {
				t.Errorf("got wrong decision %v, want %v, when proccessing %v", down, test.isDownloadable, test.URL)
			}

			if tp != test.ytType {
				t.Errorf("got wrong type %v, want %v, when proccessing %v", tp, test.ytType, test.URL)
			}

		})
	}

}
