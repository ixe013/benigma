# Enigma HTTP API

It is assumed the Enigma backend is mounted at the `enigma` path in Vault. If you mount the secret backends at any other location, update your API calls accordingly.

* [List available Enigma models](#list-available-enigma-models)
* [Create custom model](#create-custom-model)
* [Delete custom model](#delete-custom-model)

## List available Enigma models

This endpoint lists the preconfigured Enigma models.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `LIST`   | `/enigma/models`             | `200 application/json` |

### Sample request

```
curl \
    --header "X-Vault-Token: $(vault print token)" \
    --request LIST \
    $VAULT_ADDR/v1/enigma/models
```

```
vault list enigma/models
```


### Sample response

```json
[
    "I", "M3", "M4", "IXE013"
]
```

## Create custom model

This API allows to create a custom Enigma model. Nothing much you can do with it in the current version, but stay tuned!

| Method   | Path                     | Produces           |
| :------- | :----------------------- | :----------------- |
| `POST`   | `/enigma/models`         | `204 (empty body)` |


### Parameters

- `name` `(string: <required>)` – Specifies the name of the model to create.

### Sample Payload

```json
{
  "name": "my-model"
}
```

### Sample request

```
curl \
    --header "X-Vault-Token: $(vault print token)" \
    --request POST \
    --data @payload.json \
    $VAULT_ADDR/v1/enigma/models
```

```
vault write enigma/models name=my-model
```

## Delete custom model

This endpoint returns information about a named GPG key.

| Method   | Path                         | Produces               |
| :------- | :--------------------------- | :--------------------- |
| `DELETE` | `/enigma/models/:name`       | `204 (empty body)`     |

### Parameters

- `name` `(string: <required>)` – Specifies the name of the custom model to delete. This is specified as part of the URL.

### Sample request

```
$ curl \
    --header "X-Vault-Token: $(vault print token)" \
    --request DELETE
    $VAULT_ADDR/v1/enigma/models/my-model
```

```
vault delete enigma/models/my-model
```

