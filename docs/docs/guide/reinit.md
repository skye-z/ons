# Initialization

## NAS Service

After installing the NAS service, you need to:

1. Register the device with the central control server; after registration, the NAS will receive a `NAT.ID`.
2. Enable "Connection Password"; the system will generate a secure random password for you.
3. Enable "Auto Start".

## Obsidian Plugin

After installing the Obsidian plugin, in the settings you need to:

1. Enter the "Unique Identifier," which is the `NAT.ID` obtained when registering the NAS with the central control server.
2. Enter the "Connection Password."

## Central Control Service

After installing the central control service, you need to edit the `config.ini`:

1. Fill in your `Github OAuth2` credentials under the `github` section.
2. Access the central control service; the first visitor will become the admin.
3. Modify the `register` setting to decide whether to allow registrations.
