# PROTOCOLE NÉON — FPS rétro

FPS rétro à la Doom (moteur raycasting maison, style néon), écrit en Go avec
[Ebitengine](https://ebitengine.org). Cinq niveaux, 4 types d'ennemis,
3 armes, munitions, portails de sortie. Aucune dépendance d'exécution :
l'`.exe` est autonome.

## Jouer (Windows)

Double-cliquez sur **`protocole-neon.exe`** : c'est tout.

| Touche | Action |
|---|---|
| ZQSD / WASD / flèches | Se déplacer |
| Souris | Viser (curseur capturé en jeu) |
| Clic gauche | Tirer (plasma : maintenu pour tir automatique) |
| Maj (Shift) | Sprint |
| 1 / 2 / 3 ou molette | Changer d'arme |
| M | Minimap (coin haut-droit) |
| Tab | Grande carte au centre |
| Échap | Pause |
| Entrée / Espace | Commencer / recommencer |

## Armes

Débloquées au fil de la progression :

1. **Pistolet** — précis, tir semi-auto, 26 dégâts
2. **Fusil à pompe** — 6 plombs en cône, dévastateur de près, débloqué au niveau 2
3. **Fusil à plasma** — tir automatique rapide, énergie cyan, débloqué au niveau 3

Chacune a ses munitions ; les caisses `a` rechargent tout.

## Ennemis

| Ennemi | Points | Comportement |
|---|---|---|
| Imp (violet) | 100 | Poursuit et frappe |
| Spectre (rouge) | 150 | Rapide et fragile |
| Golem (vert) | 250 | Lent, très blindé, gros dégâts |
| Cracheur (vert fluo) | 200 | Garde ses distances et crache de l'acide |

## Objectif

Survivez, ramassez les **croix rouges** (+25 santé) et les **caisses de
munitions**, puis atteignez le **portail cyan** pour passer au niveau
suivant. Chaque niveau a son atmosphère : brouillard et palette de
plafond/sol différents (briques, complexe pierre/tech, cœur technology,
crypte violette, forteresse ocre).

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
go test ./...   # tests : cartes, raycasting, collisions, armes, simulation
go vet ./...
```

Structure :
- `main.go` — point d'entrée et constantes
- `levels.go` — les 5 cartes (texte ASCII)
- `ray.go` — raycasting DDA, ligne de vue, collisions
- `world.go` — état du monde, armes, ennemis, projectiles, pickups
- `render.go` — rendu murs + brouillard par distance
- `billboards.go` — sprites billboard avec z-buffer
- `assets.go` — textures et sprites générés procéduralement
- `palettes.go` — atmosphère par niveau (brouillard, plafond, sol)
- `minimap.go` — minimap et grande carte
- `hud.go` — armes, HUD, pause et écrans (titre, mort, victoire)
