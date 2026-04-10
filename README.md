<p align="center">
  <img src="scholarisbanner.png" alt="Scholaris Logo" width="664">
</p>

<div align="center">

![Go](https://img.shields.io/badge/go-1.26.0-blue?style=flat&logo=go&logoColor=white) ![WailsV3](https://img.shields.io/badge/framework-wails-red?style=flat) ![Data](https://img.shields.io/badge/data-PostgreSQL-blue?style=flat) ![Build](https://img.shields.io/badge/build-passing-brightgreen?style=flat) ![License](https://img.shields.io/badge/license-MIT-orange?style=flat)

</div>

Scholaris v2 is a cross-platform desktop application built on **Wails v3**. It utilizes a **Go** backend and a **SvelteKit** frontend, with **PostgreSQL** serving as the primary relational database.

## Architecture

The project is structured to separate system-level logic from the presentation layer:

* **Backend (`/internal`)**: Handles PostgreSQL connection pooling, migrations, and repository patterns.
* **Frontend (`/frontend`)**: A TypeScript-based SvelteKit app. Communication with Go is handled via auto-generated IPC bindings.
* **Database**: PostgreSQL for persistent data storage.

---

## Prerequisites

* **Go**: 1.21+
* **Node.js**: 18+
* **Wails v3**: [v3.wails.io](https://v3.wails.io/) (The project currently uses v3.0.0-apha.74)
* **PostgreSQL**: 14+

---

## Database Setup

### 1. Install PostgreSQL
Download and install the version appropriate for your OS from the [official PostgreSQL site](https://www.postgresql.org/download/). 

* **Default Port**: `5432`
* **Superuser**: `postgres`

### 2. Management via DataGrip
To initialize the schema and manage data:
1.  Add a **New Data Source** -> **PostgreSQL**.
2.  Enter your local credentials (Host: `localhost`, User: `postgres`).
3.  Right-click the connection and select **New > Database**. 
4.  Name the database `scholarisdb` or any name of your choosing.
5.  Update your `config.json` in the project root with your matching DB name and credentials. (You are given an config.example.json, however the file must be named config.json)

### Note:
If you plan on using other Database Managers, you can search online for the setup and deployment for your specific software.

---

## Development

### Installation
Clone the repository and install frontend dependencies:
```bash
git clone https://github.com/codex-coderex/scholaris-v2
cd scholaris-v2/frontend
npm install
```

### Build
After you made sure that you have installed all the dependencies.
Go in to the project root and run `wails3 build`.
The executable can then be found in the \bin folder.
