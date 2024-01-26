module Main exposing (..)

import Browser
import Browser.Navigation as Navigation
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as JD exposing (..) 
import Random exposing (..)


-- MODEL

type Model
    = 
    Start
    | Load_word String
    | Loading    
    | Loaded_def Definition
    | Failed String    
    | Check  Definition String
    | ShowAnswer Definition

init : () -> (Model, Cmd Msg)
init _ =
  (Start, Http.get
        { url = "../ressources/Words.txt"
        , expect = Http.expectString GotWord
        })


-- UPDATE

type Msg
    = 
    
    GotWord (Result Http.Error String)
    | GotDef(Result Http.Error Definition)    
    | RandomInt Int
    | GuessWord Definition String
    | GiveAnswer Definition
    | Restart

update :  Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        GotWord result ->
            case result of
                Ok wordlist ->(Load_word wordlist, Random.generate RandomInt(Random.int 0 99))
                Err _ -> (Failed "Failed to get the word list. Make sure to run elm reactor in from the parent folder", Cmd.none)
       
        GotDef (Ok data) ->
            (Loaded_def data, Cmd.none)
            
        GotDef (Err _) ->
            (Failed "error GotDef", Cmd.none)
       

        RandomInt x ->
            case model of
                Load_word word ->
                 case String.split " " word of
                        [] -> (Failed "The word list is empty.", Cmd.none)
                        (a::z) -> case (List.head (List.drop x (a::z))) of
                            Just answer -> (Loading, fetchData answer)
                            Nothing -> (Failed "Failed to pick a random word.", Cmd.none)
                Start -> (Failed "", Cmd.none)
                Loading ->(Loading, Cmd.none)
                Loaded_def _->   (Failed "", Cmd.none) 
                Failed _->(Failed "", Cmd.none)                
                Check _ _ ->(Failed "", Cmd.none)
                ShowAnswer _ -> (Failed "", Cmd.none)
        

        GuessWord definitions myguess ->
            (Check definitions myguess, Cmd.none)

        GiveAnswer definitions ->
            (ShowAnswer definitions, Cmd.none)

        Restart ->
            (Loading, Navigation.reload)

-- VIEW

view : Model -> Html Msg
view model = 
    div [] [
      h1 [style "font-size" "90px"][text "Guess it"]
      , button [onClick (Restart) ][text "Play again ;)"]
      , findDef model
      
    ]
    
    
    

findDef : Model -> Html Msg
findDef model = 
    case model of
            Start ->
                div [] [text "Bienvenue"]
            Loading ->
                div [] [ text "Chargement..." ]
            Load_word word->
                div [] [ text (word) ]
            


            Loaded_def definitions->            
                
                div[]
                    [ 
                    div [] [
                        input
                            [style "text-align" "center"
                            , style "font-size" "20px"
                            , style "width" "193px"
                            ,placeholder "Write here"
                            , Html.Attributes.value "", onInput (GuessWord definitions)
                            ][]
                        , button
                        [ style "font-size" "20px"
                        , style "margin-top" "3px"
                        , style "width" "200px"
                        , onClick (GiveAnswer definitions) 
                        ] [ text "Show answer" ]    
                        ]
                    , afficherDefinitions definitions.meanings
                    ]

            Check definitions myguess->                              

                if definitions.word /= myguess then
                
                    div[]
                    [ 
                    div [] [
                        input
                            [style "text-align" "center"
                            , style "font-size" "20px"
                            , style "width" "193px"
                            ,placeholder "Write here"
                            , Html.Attributes.value myguess, onInput (GuessWord definitions)
                            ][]
                        , button
                        [ style "font-size" "20px"
                        , style "margin-top" "3px"
                        , style "width" "200px"
                        , onClick (GiveAnswer definitions) 
                        ] [ text "Show answer" ]  
                        ]
                    , afficherDefinitions definitions.meanings
                    ]
                else 
                    div [][text ("Bravo, la réponse était bien " ++ definitions.word)]

            ShowAnswer definitions->
                div [][
                        input
                            [style "text-align" "center"
                            , style "font-size" "20px"
                            , style "width" "193px"
                            ,placeholder "Write here"
                            , Html.Attributes.value "", onInput (GuessWord definitions)
                            ][]
                        , button
                        [ style "font-size" "20px"
                        , style "margin-top" "3px"
                        , style "width" "200px"
                        , onClick (GiveAnswer definitions) 
                        ] [ text "Show answer" ]  
                        ,h1[] [text ("La réponse était  " ++ definitions.word)]
                        , afficherDefinitions definitions.meanings

                ]
        
                    


            Failed _->
                div [] [ text "Échec du chargement" ]




-- Fonction pour afficher une Meaning en HTML avec mise en page
afficherMeaning : Meaning -> List(Html msg)
afficherMeaning meaning =
    List.indexedMap (\index subDef ->
            if meaning.partOfSpeech /= "verb" then
                div []
                    [ p [] [ text ("• " ++ meaning.partOfSpeech ++ " définition " ++ String.fromInt (index + 1) ++ ":") ]
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


--SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions _ = Sub.none

-- MAIN


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }

--Fonctions pour decoder le site

type alias Definition =
    { word : String    
    ,meanings : List Meaning     
    }

type alias Meaning =
    { 
    partOfSpeech : String
    ,definitions : List SubDefinition    
    }

type alias SubDefinition =
    { definition : String    
    }
         
            
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
        

decodeSubDefinition : Decoder SubDefinition
decodeSubDefinition =
    JD.map SubDefinition
        (JD.field "definition" JD.string)
        
-- HTTP --

--Fonction pour créer l'url
createUrl : String -> String
createUrl word =
            "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ word
            

fetchData : String -> Cmd Msg
fetchData word =    
    Http.get
        { url = createUrl word
        , expect = Http.expectJson GotDef decodeArray
        }

