module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, map2, field, int, string)



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


type Model
  = Failure
  | Loading
  | Success Word


-- type alias Quote =
--   { quote : String
--   , source : String
--   , author : String
--   , year : Int
--   }


type alias Word =
  { word : String
  , phonetic : String
  }


init : () -> (Model, Cmd Msg)
init _ =
  (Loading, getRandomWord)



-- UPDATE


type Msg
  = MorePlease
  | GotWord (Result Http.Error Word)


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    MorePlease ->
      (Loading, getRandomWord)

    GotWord result ->
      case result of
        Ok word ->
          (Success word, Cmd.none)

        Err _ ->
          (Failure, Cmd.none)



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none



-- VIEW


view : Model -> Html Msg
view model =
  div []
    [ h2 [] [ text "Random Quotes" ]
    , viewWord model
    ]


viewWord : Model -> Html Msg
viewWord model =
  case model of
    Failure ->
      div []
        [ text "I could not load a random quote for some reason. "
        , button [ onClick MorePlease ] [ text "Try Again!" ]
        ]

    Loading ->
      text "Loading..."

    Success word ->
      div []
        [ button [ onClick MorePlease, style "display" "block" ] [ text "More Please!" ]
        , blockquote [] [ text word.word ]
        , p [ style "text-align" "right" ]
            [ text "— "
            , cite [] [ text word.word ]
            , text (" by " ++ word.phonetic)
            ]
        ]



-- HTTP


getRandomWord : Cmd Msg
getRandomWord =
  Http.get
    { url = "https://api.dictionaryapi.dev/api/v2/entries/en/computer"
    , expect = Http.expectJson GotWord wordDecoder
    }


wordDecoder : Decoder Word
wordDecoder =
  map2 Word
    (field "word" string)
    (field "phonetic" string)
