module Test2 exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style, src)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, map2, field, string, list)
import Json.Decode exposing (Error(..))
import Set exposing (..)
import Random



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
    { status : Status
    , stringListChosenWords : List String
    }


type Status
    = Failure
    | Loading
    | Success (List Word)


type alias Word =
    { word : String
    , meanings : List Meaning
    }


type alias Meaning =
    { partOfSpeech : String
    , definitions : List String
    }


wordInit : Word
wordInit = 
    { word = "hi"
    , meanings = 
        [ { partOfSpeech = "hi", definitions = ["hi"]}
        , { partOfSpeech = "hi", definitions = ["hi"]}
        ]
    }


init : () -> ( Model, Cmd Msg )
init _ =
    let
        modelInit = Model Loading []
    in
        ( modelInit, fetchListChosenWords )



-- UPDATE


type Msg
  = GotStringList (Result Http.Error String)
  | GotWord (Result Http.Error (List Word))


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        GotStringList result ->
            case result of 
                Ok stringList ->
                    let 
                        list =
                            String.split "" stringList

                        newModel = 
                            Model Loading list
            
                    in
                        ( newModel, fetchWord newModel)
                Err _ ->
                    let
                        failureModel =
                            Model (Success [wordInit]) []
                    in
                        ( failureModel, Cmd.none )
        GotWord result ->
            case result of 
                Ok word ->
                    ( { model | status = Success word }, Cmd.none)
                Err _ ->
                    ( { model | status = Failure }, Cmd.none )




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
    case model.status of
        Failure ->
            div []
                [ text "I could not load a random word for some reason. "
                ]

        Loading ->
            text "Loading..."

        Success listWord ->
            case listWord of 
                [] -> div [] [ text "There's no such word"]
                (x :: xs) ->
                    div []
                        [ blockquote [] [ text x.word ]
                        , p [ style "text-align" "right" ]
                            [ text "— "
                            , cite [] [ text x.word ]
                            -- , text (" by " ++ word.phonetic)
                            ]
                        ]


-- HTTP


-- fetchWord : Cmd Msg
-- fetchWord = 
--     Http.get 
--     { url = "https://api.dictionaryapi.dev/api/v2/entries/en/run"
--     , expect = Http.expectJson GotWord wordListDecoder
--     }


fetchWord : Model -> Cmd Msg
fetchWord model =
    let
        -- index : R
        -- index = 
        --     Random.int 0 (List.length model.stringListChosenWords - 1)
        
        -- wordRandom =
        --     case List.head (List.drop index model.stringListChosenWords) of
        --         Nothing -> "word"
        --         Just a -> a
                
        url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ "run" --wordRandom
    in
        Http.get 
            { url = url
            , expect = Http.expectJson GotWord wordListDecoder
            }


fetchListChosenWords : Cmd Msg
fetchListChosenWords =
    Http.get
    { url = "https://perso.liris.cnrs.fr/tristan.roussillon/GuessIt/thousand_words_things_explainer.txt"
    , expect = Http.expectString GotStringList
    }
    


wordListDecoder : Decoder (List Word)
wordListDecoder = 
    list wordDecoder


wordDecoder : Decoder Word
wordDecoder =
    map2 Word
        (field "word" string)
        (field "meanings" meaningListDecoder)


meaningListDecoder : Decoder (List Meaning)
meaningListDecoder =
    list meaningDecoder


meaningDecoder : Decoder Meaning
meaningDecoder = 
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" definitionsDecoder)


definitionsDecoder : Decoder (List String)
definitionsDecoder = 
    list (field "definition" string)


-- wordDecoder : Decoder String
-- wordDecoder =
--     field "word" string

-- wordListDecoder : Decoder (List String)
-- wordListDecoder = 
--     list wordDecoder


-- generateRandomIndex : Random.Generator Int
-- generateRandomIndex = 





