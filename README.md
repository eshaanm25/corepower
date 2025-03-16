# CorePower Class Reservations 🧘‍♂️

## Introduction 🌟

<img align="right" width="150" src="https://github.com/user-attachments/assets/d5375718-9e4c-43aa-86b3-9ada37939126">

CorePower is a chain of heated yoga studios. Due to its popularity, especially in larger cities like Austin, getting a reservation for classes can be challenging as spots fill up quickly.

This application was built to alleviate the frustration of missing out on classes because they're full. It automatically searches for available classes and makes reservations based on your preferences.

## How it Works 🤔

When the program is run, it:

1. Searches for available classes across chosen CorePower studios using OpenSearch
2. Finds the ideal class based on [preferences](./internal/corepower/algorithm.go)
3. Reserves a spot in the ideal class using the [CorePower API](./internal/corepower/client.go)

## Installation 🛠️

1. Clone this repository

   ```bash
   git clone https://github.com/eshaanm25/corepower.git && cd corepower
   ```

2. Run the application:

   ```bash
   make run USERNAME="your.email@example.com" PASSWORD="your_password"
   ```

## Automated Reservations ⚡

This project includes a GitHub Action that can automatically run the reservation system daily. To set this up:

1. Fork this repository
2. Go to your fork's Settings > Secrets and variables > Actions
3. Add two new secrets:
   - `COREPOWER_USERNAME`: Your CorePower username
   - `COREPOWER_PASSWORD`: Your CorePower password
4. The GitHub Action will run automatically at 12:00 AM CDT

You can also trigger the workflow manually from the Actions tab in your repository.
