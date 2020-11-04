package main

func limitServerList(all []HetznerServer, lim []string) []HetznerServer {
	var result []HetznerServer

	for _, l := range lim {
		for _, srv := range all {
			if srv.Name == l {
				result = append(result, srv)
				break
			}
		}
	}

	return result
}
