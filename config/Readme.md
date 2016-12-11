# About config

Package config is for interfacing with [log15-config](https://github.com/gernoteger/log15-config).

# Usage

Just add to your imports:

    _ "github.com/gernoteger/log15-config"
    
Be careful that the package with this import is also init'ed during execution, otherwise your configs won't resolve!
