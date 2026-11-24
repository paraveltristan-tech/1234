# PROTOCOLE NÉON — FPS rétro

FPS rétro à la Doom (moteur raycasting maison, style néon), écrit en Go avec
[Ebitengine](https://ebitengine.org). Trois niveaux, ennemis, munitions,
portails de sortie. Aucune dépendance d'exécution : l'`.exe` est autonome.

## Jouer (Windows)

Double-cliquez sur **`protocole-neon.exe`** : c'est tout.

| Touche | Action |
|---|---|
| ZQSD / WASD / flèches | Se déplacer |
| Souris | Viser (curseur capturé en jeu) |
| Clic gauche | Tirer |
| Maj (Shift) | Sprint |
| M | Minimap (coin haut-droit) |
| Tab | Grande carte au centre |
| Échap | Pause |
| Entrée / Espace | Commencer / recommencer |

Objectif : survivre, ramassez les **croix rouges** (+25 santé) et les
**boîtes jaunes** (+12 munitions), puis atteignez le **portail cyan** pour
passer au niveau suivant. Score : 100/élimination, 10/objet, 250/niveau.

## Compiler soi-même

Prérequis : [Go 1.21+](https://go.dev/dl/).

```bash
# Windows (depuis n'importe quel OS) :
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H windowsgui" -o protocole-neon.exe .

# Linux :
go build .
```

## Développement

```bash
go test ./...   # tests : cartes, raycasting, collisions, tir, simulation
go vet ./...
```

Chaque niveau a son atmosphère : brouillard et palette de plafond/sol
différents (fondations briques, complexe pierre/tech, cœur technology).

Structure :
- `main.go` — point d'entrée et constantes
- `levels.go` — les 3 cartes (texte ASCII)
- `ray.go` — raycasting DDA, ligne de vue, collisions
- `world.go` — état du monde, ennemis, pickups, tir
- `render.go` — rendu murs + sprites billboard avec z-buffer et brouillard
- `assets.go` — textures et sprites générés procéduralement
- `palettes.go` — atmosphère par niveau (brouillard, plafond, sol)
- `minimap.go` — minimap et grande carte
- `hud.go` — arme, HUD, pause et écrans (titre, mort, victoire)
