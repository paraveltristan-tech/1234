# PROTOCOLE NÉON — FPS rétro

FPS rétro à la Doom (moteur raycasting maison, style néon), écrit en Go avec
[Ebitengine](https://ebitengine.org). Trois modes : **campagne à 7 niveaux** (dont
une jungle à survoler en hélicoptère et une immense carte du désert avec VTT
et tank pilotables),
**mode survie** à vagues infinies et **mode Counter-Strike** à manches et
boutique. 5 armes (dont un sniper avec lunette), 5 types d'ennemis + 3 boss,
effets de particules (sang, impacts, acide, ions), munitions, portails de
sortie, **sons et musiques de synthèse**.
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
| E | Monter / descendre du VTT, du tank ou de l'hélicoptère |
| Clic gauche (à bord du tank) | Tirer un obus explosif (dégâts de zone) |
| Clic gauche (à bord de l'hélicoptère) | Tirer une roquette (dégâts de zone) |
| N | Couper / rétablir le son |
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

## Véhicules : VTT et TANK

Deux véhicules attendent dans le désert : le **VTT** (carré cyan sur la
minimap) et le **TANK** (carré orange). Appuie sur **E** près d'un véhicule
pour monter à bord — **E** encore pour descendre (le véhicule est déposé
derrière toi).

**VTT** — rapide :
- vitesse **×2,3** pour traverser les dunes,
- **écrasement** : rouler sur un ennemi lui inflige 75 dégâts par coup.

**TANK** — lourd et blindé :
- vitesse **×1,8**, écrasement dévastateur (**140 dégâts**),
- **CLIC GAUCHE : OBUS EXPLOSIF** — un obus part du canon et explose en
  gerbe de feu (120 dégâts au centre, zone de 2,4 cases), cadence 0,9 s,
- vue pilote dédiée : canon qui recule, frein de bouche, capot de tourelle.

À bord de l'un ou l'autre, l'arme à l'écran disparaît et le zoom lunette est
désactivé ; les véhicules restent posés où tu les abandonnes.

## Niveau 6 : la jungle

Une **canopée verte** (îlots de roche `%`, brouillard émeraude, palette
jungle) de 34 × 24 cases : 3 coffres y sont cachés et l'**hélicoptère**
attend au sud-ouest pour survoler les fourrés. Une quinzaine d'ennemis
dont plusieurs golems gardent le portail au sud-est.

## Niveau 7 : le désert

La campagne se termine sur une **carte géante de 48 × 27 cases** (le double
des niveaux précédents) : dunes, rochers épars, brume chaude et palette
sable. Le VTT et le TANK t'attendent près du spawn, 3 coffres sont
dispersés, et la sortie au sud-est est gardée par une **finale à trois
boss** : le Colosse au milieu de la carte, puis le **Dragon** et le Gardien
devant le portail. Élimine-les pour terminer la campagne.

## L'hélicoptère

Le troisième véhicule (carré **jaune** sur la minimap) parké dans la
jungle : appuie sur **E** pour monter à bord.

- **déplacement libre** : il **survole les murs** en ligne droite (borné aux
  bords de la carte), vitesse ×1,6,
- **CLIC GAUCHE : ROQUETTES** — 55 dégâts, zone de 1,3 cases, cadence
  0,3 s,
- vue pilote dédiée : verrière bulle encadrée de liserés verts,
- pas d'écrasement : on bombarde d'en haut.

## Sons et musiques

Tout l'audio est **synthétisé à la volée** (aucun fichier) : chaque arme a
son tir, la recharge claque, les coffres tintent, le dragon **rugit** avant
ses volées de feu, les explosions déflagrent. Chaque véhicule a **son
moteur** (VTT sifflant, tank diesel grave, rotor d'hélicoptère) qui monte
avec l'accélération. Deux **musiques** procédurales se relaient : thème
calme aux accords Am–F–C–G sur les écrans, basse carrée et mélodie en
combat. Touche **N** pour tout couper.

## Brûlures

Les projectiles de feu du **Dragon** enflamment le joueur : après l'impact,
tu brûles pendant **3 secondes** (−4 PV toutes les 0,5 s) avec
l'avertissement « EN FEU ! » à l'écran. Reste loin des flammes ou
enchaîne les soins.

## Passage secret (niveau 5)

Un **mur factice** (`S` sur la grande carte, rendu comme de la brique tant
qu'il n'est pas repéré) cache une pièce scellée au **nord-est du donjon**
(aile est, sous le couloir du haut) : approche-toi, le mur se révèle avec
un tintement, et tu trouves
un **coffre** ainsi qu'un **portail de skip** qui termine le niveau **sans
affronter le Gardien**. Idéal pour éviter le boss de la fin.

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
| **Dragon** (dragon executé) | 7000 | **Troisième boss** : 1100 PV, rapide (1,6) et volant, crache des **volées de 4 projectiles de feu** (18 dégâts, brûlants — infligent des **brûlures** de 3 s) ; monte à la vague 12 en survie, à la manche 8 en Counter-Strike (prime $700) et garde le portail du désert |
| **Colosse** (brute cuirassée) | 5000 | **Second boss** : 1400 PV, charge plus vite et crache des salves de 5 projectiles incandescents. Rôde au niveau 5 et au milieu du désert, frappe toutes les 10 vagues en survie et toutes les 6 manches en Counter-Strike (prime $500) |

## Objectif

Survivez, ramassez les **croix rouges** (+25 santé), les **caisses de
munitions** et les **coffres**, puis atteignez le **portail cyan** pour
passer au niveau suivant. Chaque niveau a son atmosphère : brouillard et
palette de plafond/sol différents (briques, complexe pierre/tech, cœur
technology, crypte violette, forteresse ocre, **jungle**, **désert** en finale).

Au **niveau 5**, deux boss défendent la sortie : le **Colosse**, une brute
cuirassée au noyau orange qui charge en crachant des salves de 5
projectiles incandescents, puis le **Gardien** — un démon géant au noyau
rouge — qui bloque le portail : sa barre de vie s'affiche en haut de l'écran,
il alterne coups de masse et volées de projectiles, et le portail ne s'ouvre
qu'à sa mort. Utilise le canon d'ions ! Ou cherche le **passage secret** :
un mur factice cache un portail de skip et un coffre.

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
go test ./...   # tests : cartes, raycasting, collisions, armes/ion/sniper/recharge, particules, survie, boss (gardien, colosse, dragon), CS, coffres, VTT, tank, hélicoptère, brûlures, sons, jungle, désert, simulation
go vet ./...
```

Structure :
- `main.go` — point d'entrée et constantes
- `levels.go` — les 7 cartes + arènes survie/CS + immense désert (texte ASCII)
- `ray.go` — raycasting DDA, ligne de vue, collisions
- `world.go` — état du monde, armes (dont sniper et canon d'ions), ennemis (dont dragon), véhicules (VTT, tank),
  survie, mode Counter-Strike (argent, boutique), pickups
- `particles.go` — particules (sang, impacts, acide, ions)
- `render.go` — rendu murs + brouillard par distance
- `billboards.go` — sprites billboard avec z-buffer
- `assets.go` — textures et sprites générés procéduralement
- `palettes.go` — atmosphère par niveau (brouillard, plafond, sol)
- `minimap.go` — minimap et grande carte
- `hud.go` — armes, HUD, pause et écrans (titre, mort, victoire)
- `audio.go` — synthèse procédurale des sons et musiques (PCM int16, oto)
