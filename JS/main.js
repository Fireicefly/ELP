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
  this.firstTurn = true;
  this.firstAction = false;
  this.askedTrade = false;
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

// fonction qui affiche les lettres de la main du joueur
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
  console.log("Board du " + player.name + " :");
  for (let i = 0; i < player.board.length; i++) {
    console.log("Ligne " + parseInt(i+1) + " : " + player.board[i]);
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

// Supprime le contenu du fichier de log
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
  addLog(player, "a joué le mot " + userInput);
  player.board.push(userInput);
  for (const char of userInput) {
    const index = player.hand.indexOf(char);
    if (index !== -1) {
      player.hand.splice(index, 1);
    }
  }
  draw1Letter(player);
  printLetters(player);
}

// fonction qui check si la transformation du mot est valide
function checkWordTransform(oldWord, newWord, playerHand){
  if (newWord.length < 3) return false;
  const oldWordLetters = oldWord.split('');
  const newWordLetters = newWord.split('');

  if (oldWord.length === newWord.length) return false;
  // Vérifier que chaque lettre du nouveau mot est dans la main du joueur ou était présente dans l'ancien mot
  for (let i = 0; i < newWordLetters.length; i++) {
    if (!(playerHand.includes(newWordLetters[i]) || oldWordLetters.includes(newWordLetters[i]))) {
      return false;
    }
  }
  // Vérifier que chaque lettre de l'ancien mot est présente dans le nouveau mot
  for (let i = 0; i < oldWordLetters.length; i++) {
    if (!newWordLetters.includes(oldWordLetters[i])) {
      return false;
    }
  }

  // Si toutes les conditions sont remplies, le mot est correctement transformé
  return true;
}

// fonction qui transforme un mot du plateau de jeu du joueur, ou permet de transformer un mot du plateau de jeu de l'adversaire
function transformWord(player, jarnac = false, otherPlayer = null){
    let index;
    do {
        printBoard(player);
        printLetters(player);
        index = readlineSync.question('Entrez la ligne du mot a transformer : ');
        index = parseInt(index) - 1;
        oldWord = player.board[index];
    } while (index < 0 || index >= player.board.length || oldWord === undefined);

    console.log('Vous avez choisi de transformer le mot : ' + oldWord);
    if (jarnac === false) {
        addLog(player, "a choisi de transformer le mot " + oldWord)
    } else {
        addLog(otherPlayer, "a choisi de transformer le mot " + oldWord + " de " + player.name)
    }
  
      let newWord;
    do {
        printLetters(player);
        newWord = readlineSync.question('Entrez le nouveau mot : ');
        newWord = newWord.toUpperCase();
        } while (checkWordTransform(oldWord, newWord, player.hand) === false);
    if (jarnac === false) {
        addLog(player, "a transformé le mot " + oldWord + " en " + newWord)
        player.board[index] = newWord;
    } else {
        addLog(otherPlayer, "a transformé le mot " + oldWord + " en " + newWord + " de " + player.name)
        otherPlayer.board.push(newWord);
        player.board.pop(newWord);

        const index = player.board.indexOf(newWord);
        if (index !== -1) {
            player.board.splice(index, 1);
        }
    }

    for (const char of newWord) {
    const countInNewWord = newWord.split(char).length - 1;
    const countInOldWord = oldWord.split(char).length - 1;
        if (countInNewWord > countInOldWord) {
            const excessCount = countInNewWord - countInOldWord;
            for (let i = 0; i < excessCount; i++) {
                const index = player.hand.indexOf(char);
                if (index !== -1) {
                    player.hand.splice(index, 1);
                }
            }
        }
    }
    if (jarnac === false) {
        draw1Letter(player);
        return 1;
    } else {
        return 2;
    }
}

// Fonction qui calcule le score d'un joueur
function score(player){
  let player_score = 0
  for (let i = 0; i < player.board.length; i++){
    player_score = player_score + player.board[i].length**2;
  }
  console.log(player.name,":",  player_score, "points");
  addLog(player, "a un score de " + player_score + " points")
  
}

function action_choice(player, elapsedT=0){
    let answer;
    let startTime = Date.now();
    let otherPlayer = players.filter(p => p !== player)[0];    
    do {
        answer = readlineSync.question('1 : Placer un mot   2 : Modifier un mot   3 : Passer\n');
        elapsedTime = (Date.now() + elapsedT) - startTime;
    } while ((answer !== "1" || player.hand.length < 3) && (answer !== "2" || player.board.length < 1) && answer!== "3" && answer.toLowerCase() !== "jarnac");
    if (answer ==="1"){
        return 1;
    }
    if (answer ==="2"){
        return 2;
    }
    if (answer ==="3"){
        return 3;
    }
    if ((answer.toLowerCase() === "jarnac") && !player.firstTurn && (elapsedTime) <= 3000 && otherPlayer.board.length > 0 && otherPlayer.hand.length > 0){
        return 4;
    } else if ((answer.toLowerCase() === "jarnac") && !player.firstTurn && (elapsedTime) > 3000){
        console.log("Trop tard pour Jarnac !");
        return action_choice(player, elapsedTime);
    } else if ((answer.toLowerCase() === "jarnac") && player.firstTurn){
        console.log("Impossible de faire un coup de Jarnac au premier tour");
        return action_choice(player);
    }
}

// Fonction qui effectue un coup de Jarnac
function jarnac(player){
    otherPlayer = players.filter(p => p !== player)[0];
    transformWord(otherPlayer, true, player);
}

function action(choice, player){
    if (choice === 1){
        printBoard(player);
        printLetters(player);
        addWord(player);
        return 1
    }
    if (choice === 2){
        transformWord(player);
        printBoard(player);
        return 2
    }
    if (choice === 3){
        console.log(player.name, "passe son tour");
        addLog(player, "passe son tour")
        return 3
    }
    if (choice === 4){
        addLog(player, "a fait un coup de Jarnac !")
        jarnac(player);
        return 4
    }

}

function testFinPartie(player){
    if (player.board.length === 8) {
        console.log("Fin de la partie");
        console.log("Resultats : ");
        printBoard(player1);
        score(player1);
        printBoard(player2);
        score(player2);
        return true;
    }
    return false;
}

// Echange 3 lettres
function trade3Letters(player){
    let lettersToTrade;
    let valid = false;
    let count;
    do {
        count = 0;
        printLetters(player)
        lettersToTrade = readlineSync.question('Entrez les 3 lettres que vous voulez echanger (sans espaces) : ');
        lettersToTrade = lettersToTrade.toUpperCase();
        for (const char of lettersToTrade) {
            const index = player.hand.indexOf(char);
            if (index !== -1) {
                count++;
            } else {
                break;
            }
        }
        if (count === 3) {
            for (const char of lettersToTrade) {
                const index = player.hand.indexOf(char);
                player.hand.splice(index, 1);
                letterValues[char]++;
            }
            for (let i = 0; i < 3; i++) {
                draw1Letter(player);
            }
            addLog(player, "a échangé 3 lettres")
            valid = true;
        }
    } while (!valid);
}

// Demande à un joueur s'il veut échanger 3 lettres, si oui, les échange, sinon, tire une lettre
function askTrade(player){
    let answer;
    if (player.hand.length >= 3) {
      do {
        printLetters(player);
        answer = readlineSync.question('Voulez-vous echanger 3 lettres ? (oui/non) : ');
      } while (answer.toLowerCase() !== "oui" && answer.toLowerCase() !== "non");
      if (answer.toLowerCase() === "oui") {
        trade3Letters(player);
      } else {
        draw1Letter(player);
      }
    } else {
      draw1Letter(player);
    }
}

// Création des instances de la classe Player
let player1 = new Player("Joueur 1");
let player2 = new Player("Joueur 2");
let players = [player1, player2];

let play = true
let end_player_turn = false
draw6Letters(player1);
draw6Letters(player2);
let choice;
console.log("Bienvenue au Jarnac");

function game() {
    while (play) {
      
        for (const player of players) {
          if (play) {

            do {
              let actionVal;
              if (player.firstTurn && !player.firstAction) {
                  console.log("\nAu tour du " + player.name + " :");
                  printLetters(player);
                  addWord(player);
                  choice = action_choice(player);
                  actionVal = action(choice, player);
                  player.firstAction = true;
              } else {
                  console.log("Au tour du " + player.name + " :");
                  if (!player.firstTurn) {
                    if (!player.askedTrade) {
                      askTrade(player);
                      player.askedTrade = true;
                    }
                    console.log("Vous avez 3 secondes pour faire un coup de Jarnac")
                  }
                  printBoard(player);
                  printLetters(player);
                  choice = action_choice(player);
                  actionVal = action(choice, player);
              }
              if (actionVal === 3) {
                  player.firstTurn = false;
                  player.askedTrade = false;
                  end_player_turn = true;
              }
              if (testFinPartie(player)) {
                  play = false;
                  break;
              }
          } while (end_player_turn !== true)
                  end_player_turn = false
      }
    }
  }
}


cleanLog(game);