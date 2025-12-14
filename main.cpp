#include <iostream>
#include <sstream>
#include <fstream>

#include "jsonParser.hpp"
#include "jsonValue.hpp"
#include "queryParser.hpp"

#include "hashMap.hpp"
#include "database.hpp"

using namespace std;

int main(int argc, char* argv[]) {
	DBCommand cmd = parseQuery(argc, argv);
	JSONNode document = loadCollection(cmd.database, cmd.collection);
	executeCommand(cmd, document);
	
	if (cmd.action != CommandAction::FIND) {
		saveCollection(cmd.database, cmd.collection, document);
	}
}
