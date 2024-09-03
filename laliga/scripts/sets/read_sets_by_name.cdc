import Golazos from "Golazos"

// This script returns a Set struct for the given name,
// if it exists

access(all) fun main(setName: String): Golazos.SetData {
    return Golazos.getSetDataByName(name: setName)!
}

