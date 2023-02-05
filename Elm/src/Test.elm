module Test exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style, src)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, map2, field, int, string, list)
import Json.Decode exposing (Error(..))
import Set exposing (..)
import String




-- MAIN


main : Program () Model Msg
main =
  Browser.element
    { init = init
    , update = update
    , subscriptions = subscriptions
    , view = view
    }



-- MODEL


type alias Model =
    { phrase : String
    , guesses : Set String }


-- phrase: String
-- phrase =
--   "Hi help me"


-- type alias Word =
--     word : String


-- type alias ModelWord =
--     { stringListChosenWords : String
--     , word : Word
--     }


listChosenWords : List String
listChosenWords =
    []
init : () -> (Model, Cmd Msg)
init _ =
    ( { phrase = "hi help me"
      , guesses = Set.empty
      }
    , Cmd.none 
    )



-- UPDATE


type Msg
  = Guess String
  | Restart
  | NewPhrase (Result Http.Error (List String))
  | GotText (Result Http.Error String)


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        Guess char ->
            ( { model | guesses = Set.insert char model.guesses }, Cmd.none )
        Restart ->
            ( { model | guesses = Set.empty }, fetchWord )
        NewPhrase result ->
            case result of 
                Ok phrase ->
                    case phrase of
                        [] -> 
                            ( { model | phrase = "loi me roi lan mot" }, Cmd.none )
                        (x :: xs) ->
                            ( { model | phrase = x }, Cmd.none )
                Err _ ->
                    ( { model | phrase = "loi me roi" }, Cmd.none )
        GotText result ->
            case result of
                Ok stringChosenWords ->
                    listChosenWords = (String.split " " stringChosenWords)
                Err _ ->
                    ( { model | phrase = "loi me roi" }, Cmd.none )




-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    let
        phraseHtml =
            model.phrase
                |> String.split ""
                |> List.map 
                    (\char ->
                        if char == " " then
                            " "
                        else if Set.member char model.guesses then
                            char 
                        else
                            "_"
                    )
                |> List.map 
                    (\char ->
                        span [] [ text char ]
                    )
                |> div []
        phraseSet =
            model.phrase
                |> String.split ""
                |> Set.fromList
        failuresHtml = 
            model.guesses
                |> Set.toList
                |> List.filter 
                    (\char -> not <| Set.member char phraseSet)
                |> List.map
                    (\char -> span [] [ text char ])
                |> div []
        buttonsHtml =
            "abcdefghijklmnopqrstuvwxyz"
                |> String.split ""
                |> List.map 
                    (\char ->
                        button [ onClick <| Guess char] [ text char ]
                    )
                |> div []
    in    
    div []
        [ phraseHtml
        , buttonsHtml
        , failuresHtml
        , button [ onClick Restart ] [ text "Restart" ]
        ]


-- HTTP


fetchWord : Cmd Msg
fetchWord = 
    Http.get 
    { url = "https://api.dictionaryapi.dev/api/v2/entries/en/word"
    , expect = Http.expectJson NewPhrase wordListDecoder
    }

fetchListChosenWords : Cmd Msg
fetchListChosenWords =
    Http.get
    { url = "https://perso.liris.cnrs.fr/tristan.roussillon/GuessIt/thousand_words_things_explainer.txt"
    , expect = Http.expectString GotText 
    }
wordDecoder : Decoder String
wordDecoder =
    field "word" string

wordListDecoder : Decoder (List String)
wordListDecoder = 
    list wordDecoder
