package middleware

import "net/http"

// enableCORS Le Middleware CORS est utilisé pour faire croire au navigateur que les deux serveur se situent à la meme adresse
// avant d'envoyer la requete frontend au backend, navigateur d'abord envoie une requete pré-vérification appelé OPTION (Preflight request).
// serveur Go doit intercepter cette requête, dire au navigateur "C'est bon, j'accepte", puis traiter la vraie requête POST.
// c'est ce que fait enableCORS

// EnableCORS is a middleware that allows requests from "http://localhost:5173" and supports common HTTP methods and headers.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
