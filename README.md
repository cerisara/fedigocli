# fedigocli

This is just a small test to make use of writeas/activityserve

This code creates an ActivityPub agent and serves it on a given port
It supports follow and post for federation

TODO:

- after following, need to quit and relaunch so that following becomes visible (refresh issue ?)

lorsque je cherche ce compte depuis Masto, il fait d'abord une request webfinger pour trouver l'URL du compte,
qui marche; il obtient en meme temps l'URL de outbox, il fait donc une request sur outbox, et il obtient le nb
de posts (qu'il affiche bien), et une autre url vers le 1er post ("first").
A ce moment, ca ne marche pas, car mastodon devrait aller chercher le 1er post avec cette URL "first", mais
apparemment, ca ne marche pas, il ne fait pas la requete. Pourquoi ?

