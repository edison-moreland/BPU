module \$_DFF_P_ (D, C, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, C;
    output Q;
    wire QQ;

    \$_DLATCH_P_ DLP1(.E(!C), .D(D), .Q(QQ));
    \$_DLATCH_P_ DLP2(.E(C), .D(QQ), .Q(Q));
endmodule