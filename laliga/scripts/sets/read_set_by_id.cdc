import Golazos from "Golazos"

// This script returns a Set struct for the given id,
// if it exists

access(all) fun main(id: UInt64): Golazos.SetData {
    return Golazos.getSetData(id: id)!
}

