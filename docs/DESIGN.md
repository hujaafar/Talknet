# Design system

Talknet is a place for thoughtful conversation. The interface combines an expressive opening with a calm reading surface: large imagery and type introduce the community, then readable cards and straightforward controls help people participate.

## Visual language

| Element                 | Choice              | Purpose                                    |
| ----------------------- | ------------------- | ------------------------------------------ |
| Dark surface            | `#10130e`           | Opening, navigation, and signature footer  |
| Paper surface           | `#f4f4ec`           | Discussion and account pages               |
| Card surface            | `#fffef8`           | Reading cards and forms                    |
| Ink                     | `#182016`           | Primary text on light surfaces             |
| Muted ink               | `#65705d`           | Supporting copy                            |
| Citrus                  | `#d9f06b`           | Primary actions and community note         |
| Inter                   | Local variable font | Interface and primary headlines            |
| Instrument Serif Italic | Local font          | Expressive phrases and featured discussion |

The hero conversation note uses the latest rendered discussion and links to its real detail page. It is omitted when the community is empty. Counts and topic links come from the database. Searches, topic filters, and later feed pages start directly with the discussion area.

## Motion

Scrolling translates and slightly scales the sculpture while the headline and conversation note move at separate depths. Progress and transforms are clamped at their bounds, including overscroll. A passive listener schedules at most one pending animation frame. Discussion cards reveal once when they enter the viewport.

The system’s reduced-motion preference controls the default. Visitors can explicitly enable or pause motion using the fixed control, and their choice persists locally. Small screens keep the sculpture still. Content remains visible if JavaScript fails or is disabled.

## Responsive behavior

The feed changes from topic navigation, discussion cards, and a community sidebar into a single reading column. Topic navigation scrolls horizontally on smaller screens. On phones, the opening stacks its headline, artwork, featured note, and primary action; supporting hero copy yields space to the composition. Account, profile, editor, and detail views use the same type, spacing, and focus styles.

## Editing the interface

Use shared templates for cards, reactions, navigation, and feedback. Preserve actual links and native form submissions when adding enhancements. Keep the image decorative where surrounding text already describes its purpose. Set ARIA boolean attributes to exact `true` or `false` tokens; template whitespace inside these values can change how assistive technology interprets controls.

Rebuild the application after changing embedded assets. Check the opening, feed, discussion, profile, and account flows at desktop and narrow phone widths. Verify focus visibility, readable contrast, honest empty states, and the motion preference alongside the visual treatment. The [validation record](VALIDATION.md) distinguishes completed checks from release recommendations.
