import Golazos from "Golazos"

// This script returns all the names for Set.
// These can be related to Set structs via Golazos.getSetByName() .

access(all) fun main(): [String] {
    return Golazos.getAllSetNames()
}

