module Main exposing (..)

-- import Json.Decode exposing (Decoder, map2, field, int, string, decodeString, at, index)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (..)



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


type alias Meaning =
    { partOfSpeech : String
    , definitions : String
    }


type alias Word =
    { word : String
    , meanings : Meaning
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( Loading, getRandomWord )



-- UPDATE


type Msg
    = MorePlease
    | GotWord (Result Http.Error Word)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        MorePlease ->
            ( Loading, getRandomWord )

        GotWord result ->
            case result of
                Ok word ->
                    ( Success word, Cmd.none )

                Err _ ->
                    ( Failure, Cmd.none )



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
                    -- , text (" by " ++ word.phonetic)
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
        (index 0 (field "word" string))
        meaningDecoder

meaningDecoder : Decoder Meaning
meaningDecoder =
    map2 Meaning
        (getMeanings (getStringInObjArray 0 "partOfSpeech"))
        (getMeanings (getArrayInObjArray 0 "definitions" (getStringInObjArray 0 "definition")))

getDefinition : Int -> Decoder String
getDefinition definitionIndex = getStringInObjArray definitionIndex "definitions"


getMeanings : Decoder a -> Decoder a
getMeanings decoder =
    index 0 (field "meanings" decoder)


getStringInObjArray : Int -> String -> Decoder String
getStringInObjArray objectIndex fieldName =
    index objectIndex (getStringField fieldName)

getArrayInObjArray : Int -> String -> Decoder a -> Decoder a
getArrayInObjArray objectIndex arrayName decoder = index objectIndex (field arrayName decoder)

getStringField : String -> Decoder String
getStringField name =
    field name string