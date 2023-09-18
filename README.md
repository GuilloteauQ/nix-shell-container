# Nix-shell-container

A nix shell running in a (thin) container

## Warnings

This project is just a pretext to play with containers.

Don't expect it to work perfectly :)

## Demo

![demo](demo/demo.gif)

## Use

Import this flake:

```nix
# ...
inputs.nsc.url = "github:GuilloteauQ/nix-shell-container";
# ...
```

(feel free to make the `nixpkgs` input of this flake follow yours)


Now when defining a shell use `nsc.lib.mkShell`

```nix
# ...
devShells.${system} = {
    myshell = nsc.lib.mkShell {
        # Your normal nix shell thingies
    };
};
# ...
```

By default, this function behaves like the usual `mkShell`.

You can activate the containerization by adding setting the argument `containerize` to `true` in the `mkShell`:


```nix
# ...
devShells.${system} = {
    myshell = nsc.lib.mkShell {
        containerize = true;
        # Your normal nix shell thingies
    };
};
# ...
```

Now running `nix develop .#myshell` will get you in a container with the shell environment!

(you might need to run it with `sudo` if the unprivileged user namespace are not available on your machine)

## Pass commands

It is not possible to pass commands to the shell via `nix develop .#myshell --command ...` because of the way Nix manages the `--command` flag [(see here)]https://github.com/NixOS/nix/blob/2a52ec4e928c254338a612a6b40355512298ef38/src/nix/develop.cc#L545).

The solution is to use the container shell as a `package`:


```nix
packages.${system} = {
    myenv = nsc.lib.mkShellContainer {
        # Your normal nix shell thingies
    };
};
```

To enter the container, you can run `nix run .#myenv`.

To execute commands in the container, you can run `nix run .#myenv -- COMMAND`.


## Misc

The container wrapper is written in Go, and it based on the *amazing* talk of Liz Rice: https://www.youtube.com/watch?v=8fi7uSYlOdc

