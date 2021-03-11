package com.sz.kurstracker.entities

import com.google.gson.annotations.SerializedName

data class KursData(@SerializedName("Date") val date: Long,
                    @SerializedName("BankName") val bankName: String,
                    @SerializedName("RateBuyUSD") val rateBuyUSD: Double,
                    @SerializedName("RateSellUSD") val rateSellUSD: Double,
                    @SerializedName("RateBuyEUR") val rateBuyEUR: Double,
                    @SerializedName("RateSellEUR") val rateSellEUR: Double) {
}