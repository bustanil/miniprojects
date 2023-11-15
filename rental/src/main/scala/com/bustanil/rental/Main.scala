package com.bustanil.rental

import cats.effect.{IO, IOApp}

object Main extends IOApp.Simple:
  val run = RentalServer.run[IO]
