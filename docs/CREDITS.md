# Credits and assets

Talknet began as a collaborative Go forum by [Mohamed Alasfoor](https://github.com/Mohamed-Alasfoor), [Ali Hasan](https://github.com/alihjmm), [Habib Mansoor](https://github.com/7abib04), and [Hussain Jawad](https://github.com/hujaafar). The original MIT license and commit history remain intact.

The redesign takes visual cues from Hussain’s Neo4flix and Nexora projects: strong editorial typography, a restrained accent color, large imagery, and scroll-linked motion. Talknet has its own artwork and interface; no third-party motion runtime is bundled.

## Typography

[Inter](https://github.com/rsms/inter) by Rasmus Andersson. The variable font is served locally as `static/images/inter-latin.woff2`; no external font requests or tracking are required. The upstream filename is `InterVariable.woff2`. See [the bundled SIL Open Font License](licenses/Inter-OFL.txt).

[Instrument Serif](https://github.com/google/fonts/tree/main/ofl/instrumentserif) supplies the italic editorial accent. The unmodified `InstrumentSerif-Italic.ttf` is served locally as `static/images/instrument-serif-italic.ttf`. See [its bundled SIL Open Font License](licenses/Instrument-Serif-OFL.txt).

## Conversation artwork

`static/images/conversation-art.webp` is original artwork made for this redesign using the built-in image generation tool. It is a stylized illustration, not a product photograph. The generated 1536 × 1024 PNG was converted to WebP for delivery (approximately 114 KiB). The favicon is a small source SVG; interface symbols are inline SVG.

Exact generation prompt:

> Use case: stylized-concept
> Asset type: premium editorial 3D artwork for Talknet, a conversation community website
> Primary request: two sculptural overlapping speech bubbles, one brushed chrome silver and one translucent acid-citrus yellow-green glass, floating above a dark charcoal plinth.
> Scene/backdrop: nearly black graphite studio background.
> Style/medium: exceptionally polished premium 3D editorial render with physically based reflections and subtle grain.
> Composition/framing: landscape 1536x1024, sculpture centered slightly right, composition suitable for cropping into a wide hero or square story card. Keep both speech-bubble silhouettes clearly recognizable.
> Lighting/mood: dramatic soft rim lighting, beautiful metallic highlights and luminous glass refractions, quiet sophistication.
> Color palette: brushed silver chrome, acid-citrus yellow-green, charcoal and graphite.
> Materials/textures: finely brushed chrome silver and thick translucent yellow-green glass, dark charcoal plinth.
> Constraints: no letters, numbers, text, interface chrome, logos, or watermarks.

## Sample content

The opt-in demonstration dataset contains fictional authors with `Demo` in their usernames and `example.invalid` email addresses. Their random password is discarded; these are not shared login accounts. Posts, replies, reactions, and counts are sample data. Create your own account to try the signed-in experience.
