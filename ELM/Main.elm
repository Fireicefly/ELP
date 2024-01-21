module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style)
import Html.Events exposing (onClick)
import Http
import Json.Decode as JD exposing (..) 

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
      , button [onClick GetDef][text "Reload"]
      ,findDef model
    ]
    
    

findDef model = 
    case model of
            Start ->
                div [] [text "Bienvenue"]
            Loading ->
                div [] [ text "Chargement..." ]

            


            Loaded_def definitions->                
                
                div [] [text (definitions.word)]
                
                    


            Failed ->
                div [] [ text "Échec du chargement" ]

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
    definitions : List SubDefinition
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
    list decodeDefinition
      |> andThen (\definitions ->
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
    JD.map Meaning        
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

-- HTTP
fetchData : Cmd Msg
fetchData =
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/hello"
        , expect = Http.expectJson GotDef decodeArray
        }

