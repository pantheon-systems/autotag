module github.com/pantheon-systems/autotag

go 1.12

require (
	github.com/Unknwon/com v0.0.0-20190804042917-757f69c95f3e // indirect
	github.com/gogits/git-module v0.0.0-20170404055912-2a496cad1f36
	github.com/hashicorp/go-version v0.0.0-20161031182605-e96d38404026
	github.com/jessevdk/go-flags v0.0.0-20161025193802-0648c820cd4e
	github.com/mcuadros/go-version v0.0.0-20161105183618-257f7b9a7d87 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
)

// https://github.com/pantheon-systems/autotag/issues/24
// - pulled in by: https://github.com/gogs/git-module/blob/e55accd068eac1c9803754a776c22b1aeddc4602/repo.go#L17
replace github.com/Unknwon/com v0.0.0-20190804042917-757f69c95f3e => github.com/unknwon/com v0.0.0-20190804042917-757f69c95f3e
