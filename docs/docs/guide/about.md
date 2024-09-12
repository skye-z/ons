# About ONS

ONS stands for Obsidian NAS Sync.

It is a suite dedicated to quickly and efficiently solving the file synchronization problem between Obsidian and your home NAS.

It consists of three parts:

* Obsidian Plugin: The plugin monitors changes in the vault and syncs them with the NAS.
* Central Control Server: This server handles authentication and information exchange between both ends.
* NAS Sync Server: The NAS sync service stores vault data and provides data synchronization to clients.

The **Central Control Server** can be used as `ons.betax.dev` or deployed on your own cloud server.

## Use Case

After surveying the habits of Obsidian users, ONS identified its core use case:

**1n1**: 1 NAS, multiple Obsidian clients, one online at a time.

* 1 NAS

ONS does not support multiple storage destinations.

Most Obsidian users do not require multi-location storage (less than 3%), and NAS devices often offer various sync methods with third parties.

If you need multi-location storage, use other sync services alongside, such as built-in cloud sync from your NAS.

* Multiple Obsidian Clients

ONS uses a one-way registration system. Clients only need to enter the correct NAT.ID and password, with no limit on the number of clients.

* One Online at a Time

Designed with the assumption that most users do not edit simultaneously across multiple instances, ONS lacks conflict resolution. Please ensure only one client is editing at a time.

> Note! Do not edit with multiple clients online, which may lead to content conflicts.

## Features

* Github Login: Integrated with Github for hassle-free account management.
* P2P Direct Connection: Uses WebRTC for direct P2P connection between NAS and client.
* Password Protection: NAS end supports password setting. Incorrect passwords will halt information exchange.
* Auto-Sync Changes: Obsidian client supports automatic syncing after document edits.
* Large File Chunked Transfer: Considering bandwidth limitations, `ONS` splits large files into `40KB` chunks for transfer.
