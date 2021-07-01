## gitdl

#### !! WINDOWS ONLY !!

gitdl is a tool to easily clone git repositories from the browser to the Desktop folder.

#### Usage

To build it yourself, clone the repo and run `go build`. Prebuilt executables are available in the action artifacts ([Latest build](https://nightly.link/tbhmens/gitdl/workflows/build/master/gitdl.exe.zip)).

To install, simply run `./gitdl.exe`.

To use it, go to your browser and simply replace the https:// in the repository's URL with gitdl:// and press enter. It'll automatically clone the repository and open up the folder.

For example, `https://github.com/tbhmens/gitdl` => `gitdl://github.com/tbhmens/gitdl`

Important: move the exe to a directory in which it won't be (re)moved, otherwise windows won't be able to find the gitdl executeable.