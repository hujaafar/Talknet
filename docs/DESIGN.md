# Design system

Talknet gives thoughtful conversation an expressive setting: an oversized sculptural opening leads into readable discussion cards, focused writing, and direct controls. Dark is the default, with an optional warm light reading theme.

## Visual language

| Element                 | Choice              | Purpose                                             |
| ----------------------- | ------------------- | --------------------------------------------------- |
| Dark canvas             | `#0d100e`           | Default page background                             |
| Dark surface            | `#161a17`           | Cards, forms, and the search palette                |
| Primary text            | `#f1f3ec`           | Headlines and essential controls                    |
| Muted text              | `#a2ada3`           | Supporting copy on dark surfaces                    |
| Citrus                  | `#d9f06b`           | Primary actions, highlights, and the community note |
| Light canvas            | `#f4f4ec`           | Optional paper reading theme                        |
| Inter                   | Local variable font | Interface and primary headlines                     |
| Instrument Serif Italic | Local font          | Expressive phrases                                  |

Chrome-and-citrus artwork anchors the opening. Fine orbital lines, a translucent conversation note, and restrained geometric feature cards carry the identity into the page. The three feature illustrations are native CSS geometry; they require no additional image downloads. Their cards link to actual discussions, not placeholder destinations.

The hero note uses the latest rendered discussion and is omitted for an empty community. Counts and topic links come from SQLite. Searches, topic filters, and later feed pages begin directly with the reading area.

## Motion

A staggered entrance brings in the headline. Scrolling translates and slightly scales the sculpture while the headline and live discussion note move at different depths. Progress and transforms are clamped, including overscroll. A passive listener schedules at most one pending animation frame. Discussion and feature cards reveal when they enter the viewport.

Fine-pointer desktops add small, bounded movement to the artwork stage. Geometric feature cards respond on hover. Supporting browsers animate full-page navigation with the View Transition API; ordinary navigation remains the fallback.

System reduced motion controls the default. Visitors can explicitly enable or pause motion using the fixed control, and their choice persists locally. Pause also suppresses CSS hover and navigation animation. Small screens keep the sculpture still. Essential content remains visible if JavaScript fails or is disabled.

## Useful controls

- **Search:** Ctrl/Cmd K or the header search button opens a native modal. Live results support arrow-key navigation, Enter, and Escape. The form also offers full-page search.
- **Reading:** Comfortable and Compact change card density without altering the result set. Dark/light and density preferences persist on the device.
- **Saving:** Bookmark buttons provide a private collection; a native discussion-page form provides the same action without JavaScript.
- **Writing:** Write/Preview preserves plain text, while Focus mode centers the editor. Word counts and reading estimates respond as the author types. No draft autosave or Markdown support is implied.
- **Editing:** An author can revise their discussion. Stale edits retain their text and show a conflict instead of overwriting a newer version.

## Responsive behavior

The feed changes from topic navigation, discussion cards, and a community sidebar into one reading column. Topics and feature cards scroll horizontally on smaller screens. On phones, the opening stacks typography, artwork, the live note, and its primary action. The header retains search, theme, and an account entry point. Account, profile, editor, detail, and error views share the same colors and focus treatment.

## Editing the interface

`app.css` provides the base layout; `experience.css` supplies the immersive layer. Keep theme colors tokenized and check both modes after changes. Use shared templates for cards, reactions, navigation, and feedback. Preserve actual links and native forms when adding enhancements.

Keep decorative artwork out of the accessibility tree. ARIA boolean attributes must render exact `true` or `false` tokens; template whitespace inside these values can alter their meaning. Go-template formatting can also introduce unwanted whitespace into form actions, so keep action paths explicit and check real submissions.

Rebuild after changing embedded assets. Inspect opening, feed, discussion, profile, and account flows at desktop and narrow phone widths. Check visible focus, readable contrast, useful empty states, keyboard actions, and motion preferences. The [validation record](VALIDATION.md) distinguishes completed checks from broader release recommendations.
