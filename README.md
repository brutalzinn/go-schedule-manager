# Go Schedule Manager

### IN DEVELOPMENT

The Go Schedule Manager is a program designed to help manage routines for individuals on the autism spectrum, providing constant points of reference throughout the day to easily manage changes. The application uses a web interface for scheduling events and integrating text-to-speech (TTS) capabilities.

## Features

- **Google Blockly**:  for executing blocks of code when time is trigger (future feature).
- **YouTube Integration**: Download and play audio from YouTube videos (future feature).
- **Schedule Management**: Create, update, and delete schedules.
- **Text-to-Speech Integration**: Use TTS for notifications.

- **Web Interface**: Manage schedules through a user-friendly web interface.
- **Internationalization**: Supports multiple languages using JSON translations.

## Installation

### Prerequisites

- Docker ( you dont need this but is less pain than install the voice packages for your system)
- Go
- Visual Code ( or any text editor or use vim if you cant uses mouse )

### Build and Run

1. **Clone the Repository:**
   ```sh
   $ git clone https://github.com/your-repo/go-schedule-manager.git
   $ cd go-schedule-manager
   $ docker build -t go-schedule-manager .
   $ docker run -p 8000:8000 go-schedule-manager

# Usage
Access the web interface at http://localhost:8000 to manage your schedules.

# Scheduling

## CSV
Schedules can be managed via a CSV file and converted to JSON for easy handling. The application supports time-based triggers and can play audio files or use TTS for notifications. (works like a charm to my use case)

   ```sh
      Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday
      06:30,Coffee,Coffee,Coffee,Coffee,Coffee,Coffee,Coffee
      07:00,Study,Study,Study,Study,Study,Study,Study
      08:00,Work,Work,Work,Work,Work,Work,Work
      09:00,Break,Break,Break,Break,Break,Break,Break
      10:00,Meeting,Meeting,Meeting,Meeting,Meeting,Meeting,Meeting
   ```
# Example JSON

This json is generated when you import a csv or create a shedule for especify time.

   ```sh
      [
         {
            "id": "1",
            "time": "06:30",
            "content": "Coffee",
            "useTTS": true,
         },
         {
            "id": "2",
            "time": "08:00",
            "content": "Academy",
            "useTTS": true,
         }
      ]
   ```

### Disclaimer:
   This is another hyperfocus project that i will abandon completely after my tediusus is done or my problem is solved. Feel free to contribute too.