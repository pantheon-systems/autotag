# building

on OSX brew will install libgit 0.23.0, git2go currently wants 0.22.X you can downgrade and pin with:

```
brew unlink libgit2
brew install https://raw.githubusercontent.com/Homebrew/homebrew/d6a9bb6adeb2043c5c5e9ba3a878decdefc1d240/Library/Formula/libgit2.rb
brew switch libgit2 0.22.3
brew pin libgit2
```
