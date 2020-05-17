module github.com/pantheon-systems/autotag

go 1.12

require (
	github.com/Unknwon/com v0.0.0-20190804042917-757f69c95f3e // indirect
	github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38
	github.com/alecthomas/colour v0.1.0 // indirect
	github.com/alecthomas/repr v0.0.0-20200325044227-4184120f674c // indirect
	github.com/gogits/git-module v0.0.0-20190412001904-d1773fe06f37
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mcuadros/go-version v0.0.0-20190308113854-92cdf37c5b75 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/smartystreets/assertions v1.0.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/unknwon/com v1.0.1 // indirect
)

// https://github.com/pantheon-systems/autotag/issues/24
// - pulled in by: https://github.com/gogs/git-module/blob/e55accd068eac1c9803754a776c22b1aeddc4602/repo.go#L17
replace github.com/Unknwon/com v0.0.0-20190804042917-757f69c95f3e => github.com/unknwon/com v0.0.0-20190804042917-757f69c95f3e
