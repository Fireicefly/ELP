module Test exposing (..)

import Browser
import Html exposing (Html, div, h1, input, text, button)
import Html.Attributes exposing (placeholder, type_, value)
import Html.Events exposing (onClick, onInput)


-- Model
type alias Model =
    { userInput : String
    }


-- Message
type Msg
    = UpdateUserInput String
    | GetDef


-- Initial model
initialModel : Model
initialModel =
    { userInput = ""
    }


-- Update function
update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateUserInput newText ->
            { model | userInput = newText }

        GetDef ->
            -- Handle GetDef message if needed
            model


-- View function
view : Model -> Html Msg
view model =
    div []
        [ h1 [] [text "Devinette"]
        , button [onClick GetDef] [text "Reload"]
        , findDef model
        , input
            [ type_ "text"
            , placeholder "Enter text"
            , value model.userInput  -- Use `Html.Attributes.value` directly
            , onInput UpdateUserInput
            ]
            []
        , text ("You entered: " ++ model.userInput)
        ]


findDef : Model -> Html Msg
findDef model =
    -- Implement your logic to display definitions here
    div [] [text "Definitions go here"]


-- Main program
main =
    Browser.sandbox { init = initialModel, view = view , update = update}
