# Self-Deployment

## Deployment

### Docker Container Deployment

1. Pull the `skyezhang/ons-server` image.
2. Expose ports `80` and `443`.
3. Start the container.

### Linux Deployment

```shell
bash -c "$(curl -fsSL https://betax.dev/sc/ons.sh)"
```

#### Control

```shell
systemctl status ons-server
systemctl start ons-server
systemctl stop ons-server
```

## Initial Setup

After installing the central control service, you need to edit the config.ini:

1. Enter your Github OAuth2 credentials under the github section.
2. Access the central control service; the first visitor will become the admin.
3. Modify the register setting to decide whether to enable registration.
