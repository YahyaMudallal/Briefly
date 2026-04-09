package middleware

import "net/http"

// enableCORS Le Middleware CORS est utilisé pour faire croire au navigateur que les deux serveur se situent à la meme adresse
// avant d'envoyer la requete frontend au backend, navigateur d'abord envoie une requete pré-vérification appelé OPTION (Preflight request).
// serveur Go doit intercepter cette requête, dire au navigateur "C'est bon, j'accepte", puis traiter la vraie requête POST.
// c'est ce que fait enableCORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Autoriser l'origine de votre front-end (React)
		// En développement, "*" est facile. En production, mettez "https://votresite.com"
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// 2. Autoriser les méthodes HTTP spécifiques
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		// 3. Autoriser les en-têtes spécifiques que React va envoyer
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		// 4. Gérer la requête de pré-vérification (Preflight)
		// Si le navigateur demande juste les options, on lui répond OK et on s'arrête là.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Si c'est une vraie requête (POST, GET...), on passe la main au handler principal
		next(w, r)
	}
}
