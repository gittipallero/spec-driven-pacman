# Data Model: Arcade Main Page

**Feature**: 001-main-page  
**Date**: 2026-02-03

## Overview

This feature is a static landing page with no persistent data or backend interaction. The data model is minimal, consisting only of UI configuration constants.

## Entities

### None Required

The main page does not create, read, update, or delete any persistent entities. All content is static and defined in component code.

## UI Configuration (Constants)

While not a data model in the traditional sense, the following configuration values are used:

### Theme Configuration

| Property            | Type   | Value                | Description      |
| ------------------- | ------ | -------------------- | ---------------- |
| `GAME_TITLE`        | string | "Spec-Driven-Pacman" | Main title text  |
| `START_BUTTON_TEXT` | string | "START GAME"         | Primary CTA text |

### Color Palette

Defined as CSS custom properties in `styles/variables.css`:

| Variable             | Value     | Usage                          |
| -------------------- | --------- | ------------------------------ |
| `--arcade-black`     | `#000000` | Primary background             |
| `--arcade-dark-blue` | `#0a0a23` | Secondary background           |
| `--pacman-yellow`    | `#ffff00` | Title color, primary accent    |
| `--ghost-cyan`       | `#00ffff` | Button hover, secondary accent |
| `--text-white`       | `#ffffff` | Body text                      |

## State Management

### Component State (Local Only)

| Component     | State       | Type    | Description             |
| ------------- | ----------- | ------- | ----------------------- |
| `StartButton` | `isHovered` | boolean | Hover effect trigger    |
| `StartButton` | `isPressed` | boolean | Press animation trigger |

No global state management (Context, Zustand) required for this feature.

## Future Considerations

When game features are added, the following entities will likely be needed:

- **Player**: User identity and session
- **GameState**: Current game status, score, lives
- **HighScore**: Persisted score records

These are out of scope for the main page feature.
