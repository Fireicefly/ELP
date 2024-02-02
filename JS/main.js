const readlineSync = require('readline-sync');
const fs = require('fs');
const logFile = 'log.txt';

// Définir la liste des lettres et leurs valeurs
let letterValues = {
  'A': 14, 'B': 4, 'C': 7, 'D': 5, 'E': 19, 'F': 2, 'G': 4, 'H': 2, 'I': 11,
  'J': 1, 'K': 1, 'L': 6, 'M': 5, 'N': 9, 'O': 8, 'P': 4, 'Q': 1, 'R': 10,
  'S': 7, 'T': 9, 'U': 8, 'V': 2, 'W': 1, 'X': 1, 'Y': 1, 'Z': 2
};

// Définition de la classe Player
function Player(name){
  this.name = name;
  this.hand = [];
  this.board = [];
  this.justPlayedWords = [];
}

// Création des instances de la classe Player
let player1 = new Player("Joueur 1");
let player2 = new Player("Joueur 2");
let players = [player1, player2];


// fonction pile ou face (retourne "player1" ou "player2")
function flipCoin() {
  return Math.random() < 0.5 ? "Joueur 1" : "Joueur 2";
}

// fonction qui tire 6 lettres au hasard et les ajoute à la main du joueur
function draw6Letters(player) {
  for (let i = 0; i < 6; i++) {
    draw1Letter(player);
  }
}

// fonction qui tire une lettre au hasard et l'ajoute à la main du joueur
function draw1Letter(player) {
  availableLetters = Object.keys(letterValues);
  let letter;
  do {
    const randomIndex = Math.floor(Math.random() * availableLetters.length);
    letter = availableLetters[randomIndex];
  } while (letterValues[letter] === 0);
  player.hand.push(letter);
  letterValues[letter]--;
}

function printLetters(player){
    lettres = player.hand.join(' ');
    console.log("Lettres : " + lettres);
}

// fonction qui check si le mot est valide
function checkWord(word, playerHand){
  if (word.length < 3) return false;
  for (let i = 0; i < word.length; i++) {
    if (!playerHand.includes(word[i])) {
      return false;
    }
  }
  return true;
}

// fonction qui affiche le plateau de jeu du joueur
function printBoard(player){
  console.log("Board du " + player.name + " : \n");
  for (let i = 0; i < player.board.length; i++) {
    console.log("Ligne " + parseInt(i+1) + " : " + player.board[i]+ "\n");
  }
}

// Ajoute une ligne au fichier de log, si le fichier n'existe pas, il est créé
// Utilise
function addLog(player, line) {
    fs.appendFileSync(logFile, player.name + " " + line + '\n', (err) => {
        if (err) {
            console.error("Une erreur s'est produite lors de l'écriture dans le fichier :", err);
        }
    });
}

function cleanLog(callback) {
    fs.writeFile(logFile, '', (err) => {
        if (err) {
            console.error("Une erreur s'est produite lors de la suppression du contenu du fichier :", err);
            return;
        }
        callback();
    });
}

// fonction qui ajoute un mot au plateau de jeu du joueur
function addWord(player){
  let userInput;
  do{
    userInput = readlineSync.question('Entrez un mot : ');
    userInput =userInput.toUpperCase();
  } while (checkWord(userInput, player.hand) === false);
  console.log('Vous avez saisi : ' + userInput);
  addLog(player, "a joué le mot " + userInput);
  player.board.push(userInput);
  player.justPlayedWords.push(userInput)
  player.hand = player.hand.filter(char => !userInput.split('').includes(char));
}

// fonction qui check si la transformation du mot est valide
function checkWordTransform(oldWord, newWord, playerHand){
  if (newWord.length < 3) return false;
  const oldWordLetters = oldWord.split('');
  const newWordLetters = newWord.split('');
  // Vérifier que chaque lettre du nouveau mot est dans la main du joueur et était présente dans l'ancien mot
  // for (let i = 0; i < newWordLetters.length; i++) {
  //   if (!playerHand.includes(newWordLetters[i]) || !oldWordLetters.includes(newWordLetters[i])) {
  //     return false;
  //   }
  // }
  // Vérifier que chaque lettre de l'ancien mot est présente dans le nouveau mot
  for (let i = 0; i < oldWordLetters.length; i++) {
    if (!newWordLetters.includes(oldWordLetters[i])) {
      return false;
    }
  }

  // Si toutes les conditions sont remplies, le mot est correctement transformé
  return true;
}

// fonction qui transforme un mot du plateau de jeu du joueur
function transformWord(player){
  let index;
  do {
    index = readlineSync.question('Entrez la ligne du mot a transformer : ');
    index = parseInt(index) - 1;
    oldWord = player.board[index];
  } while (index < 0 || index >= player.board.length || oldWord === undefined);
  
  console.log('Vous avez choisi de transformer le mot : ' + oldWord);
  addLog(player, "a choisi de transformer le mot " + oldWord)
  
  let newWord;
    do {
        printLetters(player);
        newWord = readlineSync.question('Entrez le nouveau mot : ');
        newWord = newWord.toUpperCase();
    } while (checkWordTransform(oldWord, newWord, player.hand) === false);
  console.log('Vous avez saisi : ' + newWord);
  addLog(player, "a transformé le mot " + oldWord + " en " + newWord)
  player.board[index] = newWord;
  player.justPlayedWords.push(newWord)
}

function scoreWord(word){
  return word.length**2;
}

function end_turn(){
    let answer;
    do {
        answer = readlineSync.question('Avez-vous terminé votre tour ?');       
        
    } while (answer !== "oui" && answer !== "non");
    if (answer ==="oui"){
        return true;
    }
    return false;
}

function action_choice(){
    let answer;
    do {
        answer = readlineSync.question('1 : Placer un mot   2 : Modifier un mot   3 : Passer\n');
        
    } while (answer !== "1" && answer !== "2" && answer!== "3");
    if (answer ==="1"){
        return 1;
    }
    if (answer ==="2"){
        return 2;
    }
    if (answer ==="3"){
        return 3;
    }
}

function action(choice, player){
    if (choice === 1){
        printBoard(player1);
        printLetters(player);
        addWord(player);
        printBoard(player1);
        return 1
    }
    if (choice === 2){
        printBoard(player1);
        transformWord(player);
        printBoard(player1);
        return 2
    }
    else{
        console.log(player.name, "passe son tour");
        addLog(player, "passe son tour")
        return 3
    }

}


let i = 0
let end_player_turn = false
draw6Letters(player1);
draw6Letters(player2);
let init = true
let choice = 0


function game() {
    while (i !== 1) {
        console.log("Bienvenue au Jarnac");
        do {
            actionVal = 0;
            console.log("Au tour du joueur 1");
            printLetters(player1);
            if (init === true) {
                addWord(player1);
                init = false;
            } else {
                choice = action_choice();
                actionVal = action(choice, player1);
            }

            if (actionVal !== 3) {
                end_player_turn = end_turn();
                if (end_player_turn === true) {
                    addLog(player1, "a terminé son tour")
                }
            } else {
                end_player_turn = true
            }


        } while (end_player_turn !== true)
        end_player_turn = false
        do {
            console.log("Au tour du joueur 2");
            printLetters(player2);
            addWord(player2);
            printBoard(player2);
            transformWord(player2);
            printBoard(player2);
            end_player_turn = end_turn()
            if (end_player_turn === true) {
                addLog(player2, "a terminé son tour")
            }
        } while (end_player_turn !== true)

    }
}

cleanLog(game);


