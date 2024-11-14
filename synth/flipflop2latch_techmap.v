// map flip flops into dlatch based implementations

module \$_DLATCH_P_ (D, E, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, E;
    output Q;

    DLATCH DLP(.E(E), .D(D), .Q(Q));
endmodule

module \$_DFF_P_ (D, C, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, C;
    output Q;
    wire QQ;

    DLATCH DLP1(.E(~C), .D(D), .Q(QQ));
    DLATCH DLP2(.E(C), .D(QQ), .Q(Q));
endmodule

module \$_DFF_PP0_ (D, C, R, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, C, R;
    output reg Q;
    wire QQ;

    DLATCH DLP1(.E(~C), .D(D), .Q(QQ));
    DLATCH DLP2(.E(C|R), .D(QQ&(~R)), .Q(Q));
endmodule