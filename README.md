# vrd

## What is vrd ?

The role of vuerd is to generate code for ORMs based on the [ERD Editor](https://marketplace.visualstudio.com/items?itemName=dineug.vuerd-vscode) vscode extension.


## Docs:

Download the binary and add it to your path, and from your command line

```console
vrd
```

that will initialize a config file `config.vrd.json`

```json
{
  "input": "db.vuerd.json",
  "output": "./",
  "ent": {
    "package": "app",
    "graphql": true,
    "echo": true,
    "auth": true,
    "privacy": true
  }
}

```

## Supported ORMs

We currently support only [Ent](https://entgo.io/) but we are planning to support more int the future.

## Contributions

we don't accept contributions at the moment, this is meant for personal use, but you are welcome to use it.
