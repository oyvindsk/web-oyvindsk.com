rm 1 2 3 4
curl "http://localhost:8080/cv/english.pdf" -o 1
curl "http://localhost:8080/cv/norwegian.pdf" -o 2
curl "http://localhost:8080/cv/cv-øyvind_gerrard_skaar-english.pdf" -o 3
curl "http://localhost:8080/cv/cv-øyvind_gerrard_skaar-norwegian.pdf" -o 4

md5sum 1 2 3 4 | sort
