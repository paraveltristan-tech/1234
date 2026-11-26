# PROTOCOLE NÉON — FPS rétro

FPS rétro à la Doom (moteur raycasting maison, style néon), écrit en Go avec
[Ebitengine](https://ebitengine.org). Trois modes : **campagne à 6 niveaux** (dont
une immense carte du désert avec un VTT pilotable),
**mode survie** à vagues infinies et **mode Counter-Strike** à manches et
boutique. 5 armes (dont un sniper avec lunette), 5 types d'ennemis, effets de
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
| R | Recharger (transfère la réserve dans le chargeur) |
| 1 / 2 / 3 / 4 / 5 ou molette | Changer d'arme |
| Clic droit (sniper équipé) | Zoom lunette (FOV réduit, visée précise) |
| B (mode Counter-Strike) | Ouvrir / fermer la boutique |
| E | Monter / descendre du VTT |
| Flèches + Entrée (titre) | Choisir le mode (campagne / survie / counter-strike) |
| M | Minimap (coin haut-droit) |
| Tab | Grande carte au centre |
| Échap | Pause |
| Entrée / Espace | Commencer / recommencer |

## Armes

Débloquées au fil de la progression :

1. **Pistolet** — précis, tir semi-auto, 26 dégâts
2. **Fusil à pompe** — 6 plombs en cône, dévastateur de près, débloqué au niveau 2
3. **Fusil à plasma** — tir automatique rapide, énergie cyan, débloqué au niveau 3
4. **Sniper** — tir unique à 95 dégâts avec une portée de 30 cases (contre 16
   pour les autres armes) : frappe net au-delà de l'horizon des couloirs.
   Maintiens le **clic droit** pour zoomer : la vue se resserre (FOV réduit),
   la sensibilité baisse et la lunette s'affiche. Débloqué au niveau 4
   (achetable $900 en Counter-Strike). Munitions : 30 max.
5. **Canon d'ions** — l'arme de boss : tir unique à 70 dégâts avec dégâts de
   zone (2,6 cases) qui touchent tous les ennemis autour du point d'impact,
   débloqué au niveau 5 (ou vague 5 en survie, achetable $1200 en
   Counter-Strike). Munitions rares (15 max).

Chacune a un **chargeur + une réserve** : tirer à vide déclenche la recharge
automatiquement, ou appuie sur **R**. Pendant la recharge l'arme descend
(animation), une barre « RECHARGE » tourne sous le réticule et tirer est
impossible jusqu'au clic final. Animation complète : l'arme descend, puis le
**chargeur se détache et tombe** (avec douilles qui éclatent à la fin de la
recharge). Les caisses `a` alimentent la réserve.

| Arme | Chargeur | Recharge |
|---|---|---|
| Pistolet | 12 | 0,9 s |
| Fusil à pompe | 8 | 1,4 s |
| Fusil à plasma | 40 | 1,2 s |
| Sniper | 5 | 1,8 s |
| Canon d'ions | 4 | 2,2 s |

Chaque arme a aussi son *skin* détaillé et son animation : culasse, pompe,
noyau d'énergie pulsant.

## Coffres au trésor

Des **coffres en bois** (`K`) sont cachés dans les niveaux (3 dans le désert,
1 déposé à chaque vague en survie et Counter-Strike). Marche dessus pour
l'ouvrir : le couvercle s'abat et révèle un butin aléatoire — gros pack de
soins, munitions pour toutes les armes, argent (Counter-Strike) ou points.
Une fois ouvert il reste sur place, visible (et repassé en gris sur la
minimap).

## Véhicule VTT

Le désert abrite un **VTT** (marqué en cyan sur la minimap). Appuie sur **E**
près du véhicule pour monter à bord :

- vitesse **×2,3** pour traverser les dunes,
- **écrasement** : rouler sur un ennemi lui inflige 75 dégâts par coup,
- tir impossible à bord (descends avec **E** pour combattre),
- pas de zoom lunette, et l'arme à l'écran disparaît.

Le VTT reste posé où tu l'abandonnes ; remonte à tout moment.

## Niveau 6 : le désert

La campagne se termine sur une **carte géante de 48 × 27 cases** (le double
des niveaux précédents) : dunes, rochers épars, brume chaude et palette
sable. Le VTT t'attend près du spawn, 3 coffres sont dispersés, et la sortie
au sud-est est gardée par une **finale à double boss** : le Colosse puis le
Gardien défendent le portail. Élimine les deux pour terminer la campagne.

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
| **Gardien** (démon géant) | 2000 | **Boss du niveau 5** : 900 PV, frappe fort et lance des volées de projectiles ; **bloque le portail** tant qu'il vit (barre de vie en haut de l'écran). Reviens toutes les 5 vagues en survie |
| **Colosse** (brute cuirassée) | 5000 | **Second boss** : 1400 PV, charge plus vite et crache des salves de 5 projectiles incandescents. Rôde au niveau 5 et au milieu du désert, frappe toutes les 10 vagues en survie et toutes les 6 manches en Counter-Strike (prime $500) |

## Objectif

Survivez, ramassez les **croix rouges** (+25 santé), les **caisses de
munitions** et les **coffres**, puis atteignez le **portail cyan** pour
passer au niveau suivant. Chaque niveau a son atmosphère : brouillard et
palette de plafond/sol différents (briques, complexe pierre/tech, cœur
technology, crypte violette, forteresse ocre, **désert** en finale).

Au **niveau 5**, deux boss défendent la sortie : le **Colosse**, une brute
cuirassée au noyau orange qui charge en crachant des salves de 5
projectiles incandescents, puis le **Gardien** — un démon géant au noyau
rouge — qui bloque le portail : sa barre de vie s'affiche en haut de l'écran,
il alterne coups de masse et volées de projectiles, et le portail ne s'ouvre
qu'à sa mort. Utilise le canon d'ions !

## Skins des armes

Les 5 armes sont dessinées à la main, pixel par pixel, avec des skins
détaillés : pompe à crosse veinée trois tons et culasse acier, plasma à
colonne d'énergie et anneaux de serpentin pulsants, sniper au corps long
liseré néon avec lunette et fenêtre de munitions, canon d'ions à double
serpentin et bandes danger jaune/noir, pistolet à glissière bicolore. Chaque
tir a son flash de bouche coloré (jaune, cyan, blanc, vert).

## Mode Counter-Strike

Depuis le titre, choisissez **COUNTER-STRIKE** : une arène sableuse en
manches infinies. Pas de ramassages au sol ici : tout passe par l'**argent**.

- Départ avec **$800** ; chaque élimination rapporte une prime
  (imp $60, spectre $80, cracheur $100, golem $150, Gardien $300,
Colosse $500) et chaque manche terminée ajoute **$300**. De plus, toutes
les 6 manches, le **Colosse** débarque.
  manche terminée ajoute **$300**.
- Appuie sur **B** pour ouvrir la **boutique** (le jeu se fige).
  **Tab** bascule entre les deux pages **MUNITIONS** et **SKINS D'ARMES** :

| Article | Prix |
|---|---|
| Fusil à pompe | $300 |
| Fusil à plasma | $700 |
| Sniper | $900 |
| Canon d'ions | $1200 |
| Soins (+40 PV) | $250 |

**Page SKINS** : achete pour chaque arme des finitions cumulables, appliquees
en temps reel a tout le dessin de l'arme :

| Finition | Prix |
|---|---|
| Chromé (acier brossé reflet bleuté) | $400 |
| Or | $700 |
| Néon (corps sombre, accents cyan) | $900 |
| Inferno (brûlé orange) | $900 |
| Toxique (vert acide) | $900 |

- Les 5 armes sont disponibles dès la première manche (à condition de payer).

## Mode survie

Depuis le titre, choisissez **SURVIE** : une arène symétrique sans portail
où les vagues arrivent sans fin. Chaque vague annonce plus d'ennemis et
d'introduit de nouveaux types (spectres dès la vague 2, cracheurs dès la
vague 3, golems dès la vague 4, Gardien géant toutes les 5 vagues), avec des
points de vie qui montent jusqu'à doubler. Fin de vague : bonus de score, +15 santé et des bonus
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
go test ./...   # tests : cartes, raycasting, collisions, armes/ion/sniper/recharge, particules, survie, boss, CS, coffres, VTT, désert, simulation
go vet ./...
```

Structure :
- `main.go` — point d'entrée et constantes
- `levels.go` — les 6 cartes + arènes survie/CS + immense désert (texte ASCII)
- `ray.go` — raycasting DDA, ligne de vue, collisions
- `world.go` — état du monde, armes (dont sniper et canon d'ions), ennemis,
  survie, mode Counter-Strike (argent, boutique), pickups
- `particles.go` — particules (sang, impacts, acide, ions)
- `render.go` — rendu murs + brouillard par distance
- `billboards.go` — sprites billboard avec z-buffer
- `assets.go` — textures et sprites générés procéduralement
- `palettes.go` — atmosphère par niveau (brouillard, plafond, sol)
- `minimap.go` — minimap et grande carte
- `hud.go` — armes, HUD, pause et écrans (titre, mort, victoire)
