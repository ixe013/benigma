# Enigma backend for Hashicorp Vault

This Hashicorp Vault secret plugin implements an Enigma machine. It's main purpose is not really about security, but providing a complete,
self-contained example that goes further than Hashicorp's sample by providing:

 - Multiple paths and operations
 - Logging
 - Persist secret state
 - Unit test framework
 - Upgradability

Read below about [installing the plugin](#download-latest-version) and the [Enigma HTTP API](#http-api).


# Download latest version

The latest version is [always available here](https://github.com/ixe013/benigma/releases/latest/download/enigma.tar.gz):

```
curl -OL https://github.com/ixe013/benigma/releases/latest/download/enigma.tar.gz
```

# Installing

## Plug-in directory

Vault will load every plugin from a single plug-in directory. It can be anywhere, but you have to make sure that only an authorized
administrator (or continuous delivery pipeline) can write to that directory. Vault's architecture makes every plugin run in a sandbox,
but a rogue plugin, if registered correctly, could read your plugin's secret

This project uses the suggested `vault/plugins` folder under your Vault installation. You can create it like this:

```
mkdir -p vault/plugins
export VAULT_PLUGINDIR=$(pwd -P)/vault/plugins
```

I will use `$VAULT_PLUGIN` from this point on in the instructions.


## Configure Vault to load plugins

Vault will only load plugins that are in a specific directory, if they are registered and the sha256 match. You must edit your
server's configuration

 1. Set the [plugin directory](https://www.vaultproject.io/docs/configuration#plugin_directory) as per [Hashicorp's documentation (can't be a symlink, beware!)](https://www.vaultproject.io/docs/internals/plugins#plugin-directory)
 1. You need to restart Vault. Sending `SIGHUP` to the process does reload the configuration, but changes to the plugin directory are not considered until you restart

If you are just trying things out with [Vault running in dev mode](https://www.vaultproject.io/docs/concepts/dev-server), start it like this.

```
vault server --dev --dev-root-token-id root --log-level trace --dev-plugin-dir=$VAULT_PLUGINDIR
```

## Install the Enigma plugin

You install this plugin in your vault instance just like you would any other plug-in. For the plugin to retain its state between
upgrades, it must:

 - Reside in the plug-in directory, so each version must have a different name 
 - Be registered with its hash, name and command to run

The first run doesn't need all of these precautions, but might has well do it all.

 1. Extract Enigma's binary to the plug-in directory: `tar xfzv enigma.tar.gz -C $VAULT_PLUGINDIR`
 1. The plugin can compute its own hash: `$VAULT_PLUGINDIR/$(tar tfz enigma.tar.gz) hash`
 1. [Register the plugin-in](https://www.vaultproject.io/api-docs/system/plugins-catalog#register-plugin) in Vault with the following command (requires sudo): ``

```
curl -i --request PUT $VAULT_ADDR/v1/sys/plugins/catalog/secret/enigma --header "X-Vault-Token: $(vault print token)" --data @- << EOF
{
  "type":"secret",
  "command":"$(tar tfz enigma.tar.gz)",
  "sha256":"$($VAULT_PLUGINDIR/$(tar tfz enigma.tar.gz) hash)"
}
EOF
```

Notice that the `command` parameter is relative to the plugin folder. The plug-in is now installed, but not enabled.


## Enable the plugin

The plugin can now be enabled like any other secret engine:

```
vault secrets enable enigma
```

# Using the plugin

The Enigma machine implements a substitution cipher, where the key is the machine model (rotors and plugboard) and the initial position
of the rotors. It is a symetric encryption scheme where the sender and receiver must share common parameters.

It means that you must create at least two instances of a given model and keep them in sync for one to be able to process text sent to
or received by the other. Another interesting point is that the Enigma, like modern symetric ciphers, does not know which cryptographic
operation it is executing. You type something on the keyboard and the lights go on. If you happen to input plaintext, the lights will be
the cipher text. If you type the cipher text on the other machine, then the lights will be the plaintext.

Typical usage between a `boat` and a `submarine` can go like this:

| # | Step | Vault command |
+---+------+---------------+
| 1 | Create the boat instances   | `vault write enigma/models/M4/instance id=boat` |
| 2 | Create the submarine instances | `vault write enigma/models/M4/instance id=submarine` |
| 3 | Type the plaintext on the boat's keyboard | `vault write enigma/instances/boat keyboard=HELLOWORLD` |
| 4 | Type the ciphertext on the submarine's keyboard | `vault write enigma/instances/submarine keyboard=MFNCZBBFZM` |
| 5 | Prepare a reply on the submarine's keyboard | `vault write enigma/instances/submarine keyboard=GOODBYE` |
| 6 | Decrypt the reply on the other machin's keyboard | `vault write enigma/instances/boat keyboard=ZNQALNO` |

## Character set

The enigma machine had only a keyboard made of letters, uppercase by convention. No digits, punctuation or space. The plugin will replace
spaces by X (before encryption) and any other non-alphabetic character will be discarded. A warning is returned when the text was altered
to fit these processing rules:

```
$ vault write enigma/instances/m4-a keyboard="Too complex for 1942!"
WARNING! The following warnings were returned from Vault:

  * Unsupported non-alphabetic characters removed from string

Key       Value
---       -----
lights    ZLQIDQGTZUMPXISE``

$ vault write enigma/instances/m4-b keyboard=ZLQIDQGTZUMPXISE
Key       Value
---       -----
lights    TOOXCOMPLEXXFORX
```

# HTTP Api

Coming soon : documentation for the full HTTP API!


# Credits

Special thanks to the original [Go Enigma implementation](https://github.com/emedvedev/enigma).

