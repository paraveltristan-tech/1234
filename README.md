# PROTOCOLE NÉON — FPS rétro

FPS rétro à la Doom (moteur raycasting maison, style néon), écrit en Go avec
[Ebitengine](https://ebitengine.org). Deux modes : **campagne à 5 niveaux**
et **mode survie** à vagues infinies. 4 types d'ennemis, 4 armes, effets de
particules (sang, impacts, acide, ions), munitions, portails de sortie.
Aucune dépendance d'exécution : l'`.exe` est autonome.

## Jouer (Windows)

Double-cliquez sur **`protocole-neon.exe`** : c'est tout.

| Touche | Action |
|---|---|
| ZQSD / WASD / flèches | Se déplacer |
| Souris | Viser (curseur capturé en jeu) |
| Clic gauche | Tirer (plasma : maintenu pour tir automatique) |
| Maj (Shift) | Sprint |
| 1 / 2 / 3 / 4 ou molette | Changer d'arme |
| Flèches + Entrée (titre) | Choisir le mode (campagne / survie) |
| M | Minimap (coin haut-droit) |
| Tab | Grande carte au centre |
| Échap | Pause |
| Entrée / Espace | Commencer / recommencer |

## Armes

Débloquées au fil de la progression :

1. **Pistolet** — précis, tir semi-auto, 26 dégâts
2. **Fusil à pompe** — 6 plombs en cône, dévastateur de près, débloqué au niveau 2
3. **Fusil à plasma** — tir automatique rapide, énergie cyan, débloqué au niveau 3
4. **Canon d'ions** — l'arme de boss : tir unique à 70 dégâts avec dégâts de
   zone (2,6 cases) qui touchent tous les ennemis autour du point d'impact,
   débloqué au niveau 5 (ou vague 5 en survie). Munitions rares (15 max).

Chacune a ses munitions ; les caisses `a` rechargent tout.

## Effets de particules

Le moteur affiche des billboards de particules avec z-buffer : gerbes de sang
(couleur par ennemi), éclats et poussière sur les impacts muraux, éclaboussures
d'acide des cracheurs, explosion ionique, scintillement vert/jaune au
ramassage des bonus. Prises en compte : rebonds au sol, gravité, amortissement
et disparition au bout de leur vie (360 particules max).

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

## Mode survie

Depuis le titre, choisissez **SURVIE** : une arène symétrique sans portail
où les vagues arrivent sans fin. Chaque vague annonce plus d'ennemis et
d'introduit de nouveaux types (spectres dès la vague 2, cracheurs dès la
vague 3, golems dès la vague 4), avec des points de vie qui montent
jusqu'à doubler. Fin de vague : bonus de score, +15 santé et des bonus
apparaissent. La palette de l'arène change à chaque vague. Score = points
d'élimination + 100 × vague terminée.

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
go test ./...   # tests : cartes, raycasting, collisions, armes/ion, particules, survie, simulation
go vet ./...
```

Structure :
- `main.go` — point d'entrée et constantes
- `levels.go` — les 5 cartes (texte ASCII)
- `ray.go` — raycasting DDA, ligne de vue, collisions
- `world.go` — état du monde, armes (dont canon d'ions), ennemis, survie, pickups
- `particles.go` — particules (sang, impacts, acide, ions)
- `render.go` — rendu murs + brouillard par distance
- `billboards.go` — sprites billboard avec z-buffer
- `assets.go` — textures et sprites générés procéduralement
- `palettes.go` — atmosphère par niveau (brouillard, plafond, sol)
- `minimap.go` — minimap et grande carte
- `hud.go` — armes, HUD, pause et écrans (titre, mort, victoire)
