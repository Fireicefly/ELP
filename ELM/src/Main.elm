module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style)
import Html.Events exposing (onClick)
import Http
import Json.Decode as JD exposing (..) 
import Random exposing (..)
import List.Extra exposing (getAt)

-- MODEL

type Model
    = 
    Start
    | Loading    
    | Loaded_def Definition
    | Failed

init : () -> (Model, Cmd Msg)
init _ =
  (Start, Cmd.none)

-- UPDATE

type Msg
    = 
    GetDef 
    | GotDef(Result Http.Error Definition)
    

update :  Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        
            
        GetDef -> 
            (Loading, fetchData)
        GotDef (Ok data) ->
            (Loaded_def data, Cmd.none)
        GotDef (Err _) ->
            (Failed, Cmd.none)

-- VIEW

view : Model -> Html Msg
view model = 
    div [] [
      h1 [][text "Devinette"]
      , button [onClick GetDef ][text "Reload"]
      ,findDef model
    ]
    
    

findDef model = 
    case model of
            Start ->
                div [] [text "Bienvenue"]
            Loading ->
                div [] [ text "Chargement..." ]

            


            Loaded_def definitions->                
                
                afficherDefinitions definitions.meanings
                
                
                    


            Failed ->
                div [] [ text "Échec du chargement" ]




-- Fonction pour afficher une Meaning en HTML avec mise en page
afficherMeaning : Meaning -> List(Html msg)
afficherMeaning meaning =
    List.indexedMap (\index subDef ->
            if meaning.partOfSpeech == "noun" then
                div []
                    [ p [] [ text ("Définition " ++ String.fromInt (index + 1) ++ ":") ]
                    , p [] [ text ("\t" ++ subDef.definition) ]
                    ]
            else
                div [][]
                ) meaning.definitions
        
-- Fonction pour afficher une liste de Meaning en HTML avec mise en page
afficherDefinitions : List Meaning -> Html msg
afficherDefinitions defs =
    div []
        (List.concatMap afficherMeaning defs)

-- Convertir List Meaning en List (List String)
convertListMeaningToListListString : List Meaning -> List (List String)
convertListMeaningToListListString listMeaning =
    List.map (\meaning ->        
            List.map (\subDef -> subDef.definition) meaning.definitions        
    ) listMeaning

--SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions _ = Sub.none

-- MAIN

main =
    Browser.element
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }

type alias Tableau = List Definition


--Fonctions pour Decoder le site

type alias Definition =
    { word : String    
    ,meanings : List Meaning     
    }

type alias Meaning =
    { 
    partOfSpeech : String
    ,definitions : List SubDefinition
    --, synonyms : List String
    --, antonyms : List String
    }

type alias SubDefinition =
    { definition : String
    --, example : Maybe String
    --, synonyms : List String
    --, antonyms : List String
    }


--removeFirstAndLastCharacter : String -> Msg
--removeFirstAndLastCharacter jsonString =
    --let
        
        --slicedString = String.dropRight 1 (String.dropLeft 1 jsonString)
    
    
    --in GotDef (jsonString)           
            
decodeArray : Decoder Definition
decodeArray  =
    JD.list decodeDefinition
      |> JD.andThen (\definitions ->
            case definitions of
                [] ->
                    -- Gérer le cas où la liste est vide
                    JD.fail "Empty list"

                firstElement :: _ ->
                    -- Récupérer le premier élément
                    JD.succeed firstElement
        )

decodeDefinition : Decoder Definition
decodeDefinition  =
    JD.map2 Definition
        (JD.field "word" JD.string)       
        (JD.field "meanings" (JD.list decodeMeaning))
               



decodeMeaning : Decoder Meaning
decodeMeaning =
    JD.map2 Meaning
        (JD.field "partOfSpeech" JD.string)        
        (JD.field "definitions" (JD.list decodeSubDefinition))
        --(Decode.field "synonyms" (Decode.list Decode.string))
        --(Decode.field "antonyms" (Decode.list Decode.string))

decodeSubDefinition : Decoder SubDefinition
decodeSubDefinition =
    JD.map SubDefinition
        (JD.field "definition" JD.string)
        --(Decode.field "example" (Decode.maybe Decode.string))
        --(Decode.field "synonyms" (Decode.list Decode.string))
        --(Decode.field "antonyms" (Decode.list Decode.string))


--decodeDefinition =
    --JD.map Definition
    --(JD.field "meanings" JD.string (JD.list (JD.field "definitions" (JD.list (JD.field "definition" JD.string)))))
    

--decodeMeanings = 
    --JD.list decodeDefinition

--Random word

motsCourants : List String
motsCourants =
    [ "Cat", "Dog", "House", "Tree", "Book", "Car", "Water", "Sun", "Moon", "Flower"
    , "Friend", "Family", "Food", "School", "Work", "Time", "Money", "Music", "Movie", "Game"
    , "Love", "Hate", "Happy", "Sad", "Big", "Small", "Hot", "Cold", "Fast", "Slow"
    , "Color", "Red", "Blue", "Green", "Yellow", "Orange", "Purple", "Black", "White"
    , "Day", "Night", "Earth", "Sky", "Ocean", "Mountain", "City", "Country", "Cloud", "Rain"
    , "Snow", "Wind", "Fire", "Ice", "Star", "Planet", "Home", "Road", "Bridge", "Street"
    , "Park", "Game", "Team", "Player", "Goal", "Idea", "Problem", "Solution", "Question", "Answer"
    , "Dream", "Sleep", "Wake", "Health", "Disease", "Doctor", "Patient", "Friend", "Enemy", "Child"
    , "Adult", "Old", "Young", "Time", "Year", "Month", "Day", "Hour", "Minute", "Second", "Future"
    , "Past", "Present", "Nature", "Environment", "Science", "Technology", "Art", "Culture", "Language", "Idea"
    ]
-- Fonction pour choisir un mot au hasard dans la liste


-- Exemple d'utilisation
--motChoisiAuHasard : Generator String
--motChoisiAuHasard =
    --getAt genererNombreAleatoire motsCourants




-- Fonction pour générer un nombre aléatoire
--genererNombreAleatoire : Int
--genererNombreAleatoire =
    --Random.generate (Random.int 0 99)
        
        

-- HTTP --

--Fonction pour créer l'url
createUrl : String
createUrl = "https://api.dictionaryapi.dev/api/v2/entries/en/pencil"

fetchData : Cmd Msg
fetchData  =    
    Http.get
        { url = createUrl
        , expect = Http.expectJson GotDef decodeArray
        }

