module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, map4, field, int, string)

-- MODEL

type Model
    = Loading
    | Loaded String
    | Failed

init : () -> (Model, Cmd Msg)
init _ =
  (Loading, fetchData)

-- UPDATE

type Msg
    = GotDef(Result Http.Error String)

update :  Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        GotDef (Ok data) ->
            (Loaded data, Cmd.none)
        GotDef (Err _) ->
            (Failed, Cmd.none)

-- VIEW

view : Model -> Html Msg
view model =
    case model of
        Loading ->
            div [] [ text "Chargement..." ]

        Loaded definition ->
            div [] [ text definition ]

        Failed ->
            div [] [ text "Échec du chargement" ]

-- MAIN

main =
    Browser.element
        { init = init
        , update = update
        , view = view
        , subscriptions = \_ -> Sub.none
        }


-- HTTP
fetchData : Cmd Msg
fetchData =
    Http.get
        { url = "https://api.dictionaryapi.dev/api/v2/entries/en/hello"
        , expect = Http.expectString GotDef 
        }