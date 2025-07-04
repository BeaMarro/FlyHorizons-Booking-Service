# ✈️ FlyHorizons - Booking Service

This is the **Booking Service** for **FlyHorizons**, an enterprise-grade flight booking system. The service is responsible for handling seat reservations, booking management, and payment event processing.

---

## 🚀 Overview

This microservice provides the core functionality for flight bookings, seat allocation, and interactions with payment and user services via **RabbitMQ**. Built with **Go** and the **Gin** framework, it connects to a **MSSQL database hosted on Azure**.

---

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Framework**: Gin
- **Database**: Microsoft SQL Server (Azure Hosted)
- **Messaging**: RabbitMQ
- **Architecture**: Microservices

---

## 📦 Features

- 🔐 **Authentication Integration**
- 📬 **Event-driven architecture** using RabbitMQ:
  - Listens to **user events** (e.g., deletion of users)
  - Responds to **payment events**
- 💳 **Handles booking flow**
- 💺 **Manages seat availability and reservation**
- 🔄 **Uses repositories for clean persistence logic**
- ⚠️ **Graceful error handling** using a centralized error module

---
📄 License
This project is shared for educational and portfolio purposes only. Commercial use, redistribution, or modification is not allowed without explicit written permission. All rights reserved © 2025 Beatrice Marro.

👤 Author
Beatrice Marro GitHub: https://github.com/beamarro
