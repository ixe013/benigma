import hvac

vault = hvac.Client()

ciphertext = vault.secrets.kv.v1.create_or_update_secret(
    method='PUT',
    mount_point='enigma',
    path='instances/boat',
    secret={
        'keyboard':'HELLO'
    }
)

print(ciphertext["data"]["lights"])
